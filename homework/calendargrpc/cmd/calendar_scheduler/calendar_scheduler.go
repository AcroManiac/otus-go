package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

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

func main() {
	// Create cancel context
	ctx, cancel := context.WithCancel(context.Background())

	// Create database connection
	conn := database.NewDatabaseConnection(
		viper.GetString("db.user"),
		viper.GetString("db.password"),
		viper.GetString("db.host"),
		viper.GetString("db.database"),
		viper.GetInt("db.port"))
	if err := conn.Init(ctx); err != nil {
		logger.Fatal("unable to connect to database", "error", err)
	}

	// Create broker manager
	manager := broker.NewManager(
		viper.GetString("amqp.protocol"),
		viper.GetString("amqp.user"),
		viper.GetString("amqp.password"),
		viper.GetString("amqp.host"),
		viper.GetInt("amqp.port"))
	if err := manager.Open(); err != nil {
		logger.Fatal("error initializing RabbitMQ broker", "error", err)
	}

	// Set interrupt handler
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start scheduler logic in a separate goroutine
	go func() {
		scheduler := logic.NewScheduler(
			database.NewDatabaseEventsCollector(conn),
			database.NewDatabaseCleaner(
				conn,
				logic.NewRetentionPolicy(365*24*time.Hour)),
			manager.GetWriter())
		scheduleTicker := time.NewTicker(viper.GetDuration("app.scheduler_interval"))
		cleanTicker := time.NewTicker(viper.GetDuration("app.cleaner_interval"))
	OUTER:
		for {
			select {
			case <-ctx.Done():
				logger.Debug("Exit from schedule logic")
				break OUTER
			case <-scheduleTicker.C:
				if err := scheduler.Schedule(); err != nil {
					logger.Error("scheduler error", "error", err)
				}
			case <-cleanTicker.C:
				if err := scheduler.Clean(); err != nil {
					logger.Error("cleaner error", "error", err)
				}
			}
		}
	}()

	logger.Info("Application started. Press Ctrl+C to exit...")

	// Wait for user or OS interrupt
	<-done

	// Call context to stop i/o operations
	cancel()

	// Make broker graceful shutdown
	if err := manager.Close(); err != nil {
		logger.Error("failed closing RabbitMQ broker connection", "error", err)
	}

	logger.Info("Application exited properly")
}
