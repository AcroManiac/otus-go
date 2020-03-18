package main

import (
	"context"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/interfaces"
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
	application.Init("../../configs/calendar_sender.yaml")
}

func main() {
	// Create cancel context
	ctx, cancel := context.WithCancel(context.Background())

	// Create database connection
	db := database.NewDatabaseConnection(
		viper.GetString("db.user"),
		viper.GetString("db.password"),
		viper.GetString("db.host"),
		viper.GetString("db.database"),
		viper.GetInt("db.port"))
	if err := db.Init(ctx); err != nil {
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

	// Create sender logic
	sender := logic.NewSender(
		manager.GetReader(ctx),
		[]interfaces.Sender{
			logger.NewLogSender(),
			database.NewDatabaseSender(db),
		})

	// Start sender logic in a separate goroutine
	go sender.Start(ctx)

	// Set interrupt handler
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Application started. Press Ctrl+C to exit...")

	// Wait for user or OS interrupt
	<-done

	// Call context to stop i/o operations
	cancel()

	// Make broker  graceful shutdown
	if err := manager.Close(); err != nil {
		logger.Error("failed closing RabbitMQ broker connection", "error", err)
	}

	logger.Info("Application exited properly")
}
