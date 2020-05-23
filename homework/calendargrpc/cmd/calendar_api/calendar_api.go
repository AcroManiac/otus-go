package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ahamtat/otus-go/homework/calendargrpc/internal/infrastructure/application"

	"github.com/ahamtat/otus-go/homework/calendargrpc/internal/infrastructure/grpcapi"

	"github.com/ahamtat/otus-go/homework/calendargrpc/pkg/api"

	"github.com/ahamtat/otus-go/homework/calendargrpc/internal/infrastructure/database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ahamtat/otus-go/homework/calendargrpc/internal/domain/logic"
	"github.com/ahamtat/otus-go/homework/calendargrpc/internal/infrastructure/logger"
	"github.com/spf13/viper"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	// Create listener for gRPC server
	lis, err := net.Listen("tcp",
		fmt.Sprintf("%s:%d", viper.GetString("grpc.ip"), viper.GetInt("grpc.port")))
	if err != nil {
		logger.Fatal("failed to listen tcp", "error", err)
	}

	// Create a gRPC Server with gRPC interceptor
	grpcServer := grpc.NewServer()
	api.RegisterCalendarApiServer(grpcServer, grpcapi.NewCalendarApiServer(cal))

	// Register reflection service on gRPC server
	reflection.Register(grpcServer)

	// Expose the registered metrics via HTTP
	httpServer := &http.Server{
		Handler: promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}),
		Addr:    fmt.Sprintf("0.0.0.0:%d", viper.GetInt("monitor.port")),
	}

	// Set interrupt handler
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start http server for Prometheus
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Fatal("Unable to start a http server", "error", err)
		}
	}()
	logger.Info("Prometheus HTTP server started")

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
