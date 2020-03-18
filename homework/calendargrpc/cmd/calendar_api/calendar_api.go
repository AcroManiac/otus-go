package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/application"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/grpcapi"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/pkg/api"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/logic"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/logger"
	"github.com/spf13/viper"
)

func init() {
	application.Init("../../configs/calendar_api.yaml")
}

func main() {
	//// Start tracer for debug needs
	//trace.Start(os.Stderr)
	//defer trace.Stop()

	// Create cancel context
	ctx, cancel := context.WithCancel(context.Background())

	// Create database connection
	storage := database.NewDatabaseStorage(ctx,
		viper.GetString("db.user"),
		viper.GetString("db.password"),
		viper.GetString("db.host"),
		viper.GetString("db.database"),
		viper.GetInt("db.port"))

	// Create calendar
	cal := logic.NewCalendar(storage)
	logger.Info("Calendar business logic created")

	// Initialize and start gRPC server
	lis, err := net.Listen("tcp",
		fmt.Sprintf("%s:%d", viper.GetString("grpc_listen.ip"), viper.GetInt("grpc_listen.port")))
	if err != nil {
		logger.Fatal("failed to listen tcp", "error", err)
	}
	grpcServer := grpc.NewServer()
	api.RegisterCalendarApiServer(grpcServer, grpcapi.NewCalendarApiServer(cal))

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	// Set interrupt handler
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Listen gRPC server
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("error while starting gRPC server", "error", err)
		}
	}()
	logger.Info("gRPC server started. Press Ctrl+C to exit...")

	// Wait for user or OS interrupt
	<-done

	// Call context to stop i/o operations
	cancel()

	// Make gRPC server graceful shutdown
	grpcServer.GracefulStop()
	logger.Info("gRPC server stopped gracefully")

	logger.Info("Application exited properly")
}
