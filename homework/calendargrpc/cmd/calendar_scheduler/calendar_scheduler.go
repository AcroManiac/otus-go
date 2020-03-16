package main

import (
	"context"
	"flag"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/logic"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/broker"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/database"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/logger"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	// using standard library "flag" package
	flag.String("config", "../../configs/calendar_scheduler.yaml", "path to configuration flag")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	// Reading configuration from file
	configPath := viper.GetString("config") // retrieve value from viper
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Couldn't read configuration file: %s", err.Error())
	}

	// Setting log parameters
	logger.Init(viper.GetString("log.log_level"), viper.GetString("log.log_file"))
}

func main() {
	// Create cancel context
	ctx, cancel := context.WithCancel(context.Background())

	// Create database connection
	collector := database.NewDatabaseEventsCollector(ctx,
		viper.GetString("db.user"),
		viper.GetString("db.password"),
		viper.GetString("db.host"),
		viper.GetString("db.database"),
		viper.GetInt("db.port"))

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

	// Create scheduler
	scheduler := logic.NewScheduler(collector, manager.GetWriter())

	// Start scheduler logic
	ticker := time.NewTicker(10 * time.Second) //1 * time.Minute)
OUTER:
	for {
		select {
		case <-done:
			logger.Debug("Exit from ticker")
			break OUTER
		case <-ticker.C:
			if err := scheduler.Schedule(); err != nil {
				logger.Error("scheduler error", "error", err)
			}
		}
	}

	// Make broker graceful shutdown
	if err := manager.Close(); err != nil {
		logger.Error("failed closing RabbitMQ broker connection", "error", err)
	}

	// Call context to stop i/o operations
	cancel()
	logger.Info("Application exited properly")
}
