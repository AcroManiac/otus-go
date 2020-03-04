package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/api"
	"google.golang.org/grpc"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/logger"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	// using standard library "flag" package
	flag.String("config", "../../configs/calendar_api_client.yaml", "path to configuration flag")

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
	ctx, _ := context.WithCancel(context.Background())

	// Start gRPC client
	cc, err := grpc.Dial(
		fmt.Sprintf("%s:%d", viper.GetString("grpc.ip"), viper.GetInt("grpc.port")),
		grpc.WithInsecure())
	if err != nil {
		logger.Fatal("could not connect gRPC server", "error", err)
	}
	defer cc.Close()

	grpcClient := api.NewCalendarApiClient(cc)

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		logger.Error("error loading location", "error", err)
	}
	startTime, err := ptypes.TimestampProto(
		time.Date(2020, 3, 11, 20, 0, 0, 0, loc))
	if err != nil {
		logger.Error("error converting timestamp", "error", err)
	}
	_, err = grpcClient.CreateEvent(ctx, &api.CreateEventRequest{
		Title:       "Срок сдачи ДЗ",
		Description: "Срок сдачи домашнего задания №22",
		StartTime:   startTime,
		Duration:    ptypes.DurationProto(time.Hour),
	})
	if err != nil {
		logger.Error("failed sending CreateEvent", "error", err)
	}

	logger.Info("Client exited")
}
