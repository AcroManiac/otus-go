package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/interfaces"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/logic"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/application"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/broker"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/database"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/logger"
	"github.com/spf13/viper"
)

func init() {
	application.Init("../../configs/calendar_scheduler.yaml")
}

// appServices holds objects to communicate with services
type appServices struct {
	ctx       context.Context
	cancel    context.CancelFunc
	conn      *database.Connection
	manager   *broker.Manager
	scheduler interfaces.Scheduler
}

// Create scheduler logic and start in a separate goroutine
func startScheduler(app *appServices) {

	// Create cancel context
	app.ctx, app.cancel = context.WithCancel(context.Background())

	// Start connection listener
	go app.manager.ConnectionListener(app.ctx)

	// Create scheduler
	app.scheduler = logic.NewScheduler(
		app.ctx,
		database.NewDatabaseEventsCollector(app.conn),
		database.NewDatabaseCleaner(
			app.conn,
			logic.NewRetentionPolicy(viper.GetDuration("app.retention"))),
		app.manager.GetWriter(),
		viper.GetDuration("app.scheduler"),
		viper.GetDuration("app.cleaner"))
	go app.scheduler.Start()
}

func main() {

	app := &appServices{}

	// Create database connection
	app.conn = database.NewDatabaseConnection(
		viper.GetString("db.user"),
		viper.GetString("db.password"),
		viper.GetString("db.host"),
		viper.GetString("db.database"),
		viper.GetInt("db.port"))
	if err := app.conn.Init(context.Background()); err != nil {
		logger.Fatal("unable to connect to database", "error", err)
	}

	// Create broker manager
	app.manager = broker.NewManager(
		viper.GetString("amqp.protocol"),
		viper.GetString("amqp.user"),
		viper.GetString("amqp.password"),
		viper.GetString("amqp.host"),
		viper.GetInt("amqp.port"))
	if app.manager == nil {
		logger.Fatal("failed connecting to RabbitMQ")
	}
	logger.Info("RabbitMQ broker connected", "host", viper.GetString("amqp.host"))

	// Set interrupt handler
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Initialize and start scheduler
	startScheduler(app)

	logger.Info("Application started. Press Ctrl+C to exit...")

OUTER:
	for {
		select {
		// Wait for user or OS interrupt
		case <-done:
			break OUTER

		// Catch broker connection notification
		case connErr := <-app.manager.Done:
			if connErr != nil {
				// Call context to stop i/o operations and scheduler
				app.cancel()

				// Recreate broker connection and scheduler
				if err := app.manager.Reconnect(); err != nil {
					logger.Error("error reconnecting RabbitMQ", "error", err)
					break OUTER
				}
				startScheduler(app)
			}
		}
	}

	// Call context to stop i/o operations
	app.cancel()

	// Make broker graceful shutdown
	if err := app.manager.Close(); err != nil {
		logger.Error("failed closing RabbitMQ broker connection", "error", err)
	}

	logger.Info("Application exited properly")
}
