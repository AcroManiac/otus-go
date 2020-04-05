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
	application.Init("../../configs/calendar_sender.yaml")
}

// Global variables
var (
	ctx     context.Context
	cancel  context.CancelFunc
	conn    *database.Connection
	manager *broker.Manager
	sender  *logic.Sender
)

func startSender() {

	// Create cancel context
	ctx, cancel = context.WithCancel(context.Background())

	// Start connection listener
	go manager.ConnectionListener(ctx)

	// Create sender logic
	sender = logic.NewSender(
		manager.GetReader(ctx),
		[]interfaces.Sender{
			logger.NewLogSender(),
			database.NewDatabaseSender(conn),
		})

	// Start sender logic in a separate goroutine
	go sender.Start(ctx)
}

func main() {

	// Create database connection
	conn = database.NewDatabaseConnection(
		viper.GetString("db.user"),
		viper.GetString("db.password"),
		viper.GetString("db.host"),
		viper.GetString("db.database"),
		viper.GetInt("db.port"))
	if err := conn.Init(context.Background()); err != nil {
		logger.Fatal("unable to connect to database", "error", err)
	}

	// Create broker manager
	manager = broker.NewManager(
		viper.GetString("amqp.protocol"),
		viper.GetString("amqp.user"),
		viper.GetString("amqp.password"),
		viper.GetString("amqp.host"),
		viper.GetInt("amqp.port"))
	if manager == nil {
		logger.Fatal("failed connecting to RabbitMQ")
	}
	logger.Info("RabbitMQ broker connected", "host", viper.GetString("amqp.host"))

	// Set interrupt handler
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Initialize and start scheduler
	startSender()

	logger.Info("Application started. Press Ctrl+C to exit...")

OUTER:
	for {
		select {
		// Wait for user or OS interrupt
		case <-done:
			break OUTER

		// Catch broker connection notification
		case connErr := <-manager.Done:
			if connErr != nil {
				// Call context to stop i/o operations and scheduler
				cancel()

				// Recreate broker connection and scheduler
				if err := manager.Reconnect(); err != nil {
					logger.Error("error reconnecting RabbitMQ", "error", err)
					break OUTER
				}
				startSender()
			}
		}
	}

	// Call context to stop i/o operations
	cancel()

	// Make broker  graceful shutdown
	if err := manager.Close(); err != nil {
		logger.Error("failed closing RabbitMQ broker connection", "error", err)
	}

	logger.Info("Application exited properly")
}
