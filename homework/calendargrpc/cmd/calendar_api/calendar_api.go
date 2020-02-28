package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/logic"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/logger"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	// using standard library "flag" package
	flag.String("config", "../../configs/calendar_api.yaml", "path to configuration flag")

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
	//// Start tracer for debug needs
	//trace.Start(os.Stderr)
	//defer trace.Stop()

	// Create database connection

	// Create calendar
	cal := logic.NewCalendar(storage.NewStorage())
	logger.Info("Calendar was created")

	// Initialize and start gRPC server

	// Set interrupt handler
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		//// Listen gRPC server
		//if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		//	logger.Fatal("Error while starting gRPC server", "error", err)
		//}
	}()
	logger.Info("gRPC server started")

	<-done
	logger.Info("gRPC server stopped")

	/*ctx*/
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//// Make gRPC server graceful shutdown
	//if err := srv.Shutdown(ctx); err != nil {
	//	logger.Fatal("gRPC server shutdown failed", "error", err)
	//}
	logger.Info("gRPC server exited properly")
}
