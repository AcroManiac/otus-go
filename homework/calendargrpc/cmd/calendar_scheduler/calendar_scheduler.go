package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

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
	if manager == nil {
		logger.Fatal("failed connecting to RabbitMQ")
	}
	logger.Info("RabbitMQ broker connected", "host", viper.GetString("amqp.host"))

	// Set interrupt handler
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Create scheduler logic and start in a separate goroutine
	scheduler := logic.NewScheduler(
		ctx,
		database.NewDatabaseEventsCollector(conn),
		database.NewDatabaseCleaner(
			conn,
			logic.NewRetentionPolicy(viper.GetDuration("app.retention"))),
		manager.GetWriter(),
		viper.GetDuration("app.scheduler"),
		viper.GetDuration("app.cleaner"))
	go scheduler.Start()

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
