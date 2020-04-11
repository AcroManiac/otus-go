package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/application"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/pkg/api"

	"github.com/golang/protobuf/ptypes"

	"google.golang.org/grpc"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/logger"
	"github.com/spf13/viper"
)

type idsArray struct {
	mx  *sync.Mutex
	ids []string
}

func init() {
	application.Init("../../configs/calendar_api_client.yaml")
	rand.Seed(time.Now().UnixNano())
}

func main() {

	createdEvs := &idsArray{
		mx:  new(sync.Mutex),
		ids: make([]string, 0, 1000),
	}

	// Create cancel context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set interrupt handler
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start gRPC client
	cc, err := grpc.Dial(
		fmt.Sprintf("%s:%d", viper.GetString("grpc.ip"), viper.GetInt("grpc.port")),
		grpc.WithInsecure())
	if err != nil {
		logger.Fatal("could not connect gRPC server", "error", err)
	}
	defer cc.Close()

	grpcClient := api.NewCalendarApiClient(cc)

	// Create and send create events to gRPC server
	go func() {
		var eventIndex int
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.NewTimer(time.Duration(rand.Intn(5)) * time.Second).C:
				startTime, err := ptypes.TimestampProto(time.Now().Add(time.Minute))
				if err != nil {
					logger.Error("error converting timestamp", "error", err)
				}

				// Send create event request to gRPC server
				createResponse, err := grpcClient.CreateEvent(ctx, &api.CreateEventRequest{
					Title:       fmt.Sprintf("Event #%d", eventIndex),
					Description: "Data for testing microservices",
					Owner:       "Artem",
					StartTime:   startTime,
					Duration:    ptypes.DurationProto(time.Minute),
					Notify:      ptypes.DurationProto(time.Second),
				})
				if err != nil {
					logger.Error("failed sending CreateEvent request", "error", err)
					continue
				}
				eventIndex++

				// Get new event id from gRPC response
				respEvent := createResponse.GetEvent()
				if respEvent == nil {
					logger.Error("response returned no event")
					continue
				}
				logger.Debug("Event created in calendar")

				// Save event id
				createdEvs.mx.Lock()
				createdEvs.ids = append(createdEvs.ids, respEvent.GetId())
				createdEvs.mx.Unlock()
			}
		}
	}()

	// Get events from gRPC server
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.NewTimer(time.Duration(rand.Intn(3)) * time.Minute).C:
				searchTime, err := ptypes.TimestampProto(time.Now())
				if err != nil {
					logger.Error("error converting timestamp", "error", err)
				}
				searchResponse, err := grpcClient.GetEvents(ctx, &api.GetEventsRequest{
					Period:    api.TimePeriod_TIME_DAY,
					StartTime: searchTime,
				})
				if err != nil {
					logger.Error("failed sending GetEvents request", "error", err)
				}

				// Get events from response
				events := searchResponse.GetEvents()
				if events == nil {
					logger.Info("response returned no event")
				}
				logger.Info(fmt.Sprintf("GetEvents returned %d events", len(events)))
			}
		}
	}()

	// Delete events from calendar
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.NewTimer(time.Duration(2) * time.Minute).C:

				// Copy ids and clear array
				createdEvs.mx.Lock()
				copiedIds := make([]string, len(createdEvs.ids))
				copy(copiedIds, createdEvs.ids)
				createdEvs.ids = make([]string, 0, 1000)
				createdEvs.mx.Unlock()

				logger.Info(fmt.Sprintf("Deleting %d events", len(copiedIds)))
				for _, id := range copiedIds {
					// Delete event through gRPC request
					_, err = grpcClient.DeleteEvent(ctx, &api.DeleteEventRequest{
						Id: id,
					})
					if err != nil {
						logger.Error("failed sending DeleteEvent request", "error", err)
					}
				}
			}
		}
	}()

	// Wait for user or OS interrupt
	<-done

	// Call context to stop i/o operations
	cancel()

	logger.Info("Client exited")
}

func mainPrevious() {

	// Create cancel context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
		logger.Fatal("error converting timestamp", "error", err)
	}

	// Send create event request to gRPC server
	createResponse, err := grpcClient.CreateEvent(ctx, &api.CreateEventRequest{
		Title:       "Срок сдачи ДЗ",
		Description: "Срок сдачи домашнего задания №22",
		Owner:       "Артём",
		StartTime:   startTime,
		Duration:    ptypes.DurationProto(time.Hour),
	})
	if err != nil {
		logger.Fatal("failed sending CreateEvent request", "error", err)
	}

	// Get new event id from gRPC response
	respEvent := createResponse.GetEvent()
	if respEvent == nil {
		logger.Fatal("response returned no event")
	}
	logger.Debug("Event created in calendar")

	id := respEvent.GetId()

	// Send edit event request
	newTime, err := ptypes.TimestampProto(
		time.Date(2020, 3, 11, 12, 0, 0, 0, loc))
	if err != nil {
		logger.Fatal("error converting timestamp", "error", err)
	}
	_, err = grpcClient.EditEvent(ctx, &api.EditEventRequest{
		Id: id,
		Event: &api.Event{
			Id:          id,
			Title:       respEvent.GetTitle(),
			Description: respEvent.GetDescription() + " (изменено)",
			Owner:       respEvent.GetOwner(),
			StartTime:   newTime,                             // change start time
			Duration:    ptypes.DurationProto(2 * time.Hour), // and duration
			Notify:      respEvent.GetNotify(),
		},
	})
	if err != nil {
		logger.Fatal("failed sending EditEvent request", "error", err)
	}
	logger.Debug("Event edited in calendar")

	// Get events from gRPC server
	searchTime, err := ptypes.TimestampProto(
		time.Date(2020, 3, 11, 10, 0, 0, 0, loc))
	if err != nil {
		logger.Fatal("error converting timestamp", "error", err)
	}
	searchResponse, err := grpcClient.GetEvents(ctx, &api.GetEventsRequest{
		Period:    api.TimePeriod_TIME_DAY,
		StartTime: searchTime,
	})
	if err != nil {
		logger.Fatal("failed sending GetEvents request", "error", err)
	}

	// Get events from response
	events := searchResponse.GetEvents()
	if events == nil {
		logger.Info("response returned no event")
	}
	for _, ev := range events {
		logger.Info("returned event",
			"id", ev.Id,
			"title", ev.Title,
			"description", ev.Description,
			"owner", ev.Owner,
			"start_time", ev.StartTime,
			"duration", ev.Duration,
			"notify", ev.Notify)
	}

	// Delete event through gRPC request
	_, err = grpcClient.DeleteEvent(ctx, &api.DeleteEventRequest{
		Id: id,
	})
	if err != nil {
		logger.Fatal("failed sending DeleteEvent request", "error", err)
	}
	logger.Debug("Event deleted successfully")

	logger.Info("Client exited")
}
