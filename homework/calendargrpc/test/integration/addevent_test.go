package main

import (
	"context"
	"encoding/json"
	"time"

	"google.golang.org/grpc"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/entities"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/pkg/api"
	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
)

var (
	err            error
	clientConn     *grpc.ClientConn
	grpcClient     api.CalendarApiClient
	eventData      *api.CreateEventRequest
	createResponse *api.CreateEventResponse
)

func connectionToCalendarAPIOn(arg1 string) error {

	// Start gRPC client
	clientConn, err = grpc.Dial(arg1, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "could not connect gRPC server")
	}

	grpcClient = api.NewCalendarApiClient(clientConn)
	if grpcClient == nil {
		return errors.New("failed creating calendar API client")
	}

	return nil
}

func thereIsTheEvent(arg1 *messages.PickleStepArgument_PickleDocString) error {
	if arg1 == nil {
		return errors.New("argument is invalid")
	}

	event := &entities.Event{}
	if err = json.Unmarshal([]byte(arg1.Content), event); err != nil {
		return errors.Wrap(err, "couldn't parse JSON object")
	}

	startTime, err := ptypes.TimestampProto(event.StartTime)
	if err != nil {
		return errors.Wrap(err, "error converting timestamp")
	}

	eventData = &api.CreateEventRequest{
		Title:       event.Title,
		Description: *event.Description,
		Owner:       event.Owner,
		StartTime:   startTime,
		Duration:    ptypes.DurationProto(time.Minute),
		Notify:      ptypes.DurationProto(time.Second),
	}
	return nil
}

func iSendAddEventRequest() error {

	// Create cancel context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Send create event request to gRPC server
	createResponse, err = grpcClient.CreateEvent(ctx, eventData)
	if err != nil {
		return errors.Wrap(err, "failed sending CreateEvent request")
	}

	return nil
}

func responseShouldHaveEvent() error {

	// Get new event id from gRPC response
	respEvent := createResponse.GetEvent()
	if respEvent == nil {
		return errors.New("response returned no event")
	}
	if len(respEvent.GetId()) == 0 {
		return errors.New("no valid event ID returned")
	}
	return nil
}

func closeClient(*messages.Pickle, error) {
	if clientConn != nil {
		_ = clientConn.Close()
	}
}

func FeatureContext(s *godog.Suite) {

	s.Step(`^Connection to Calendar API on "([^"]*)"$`, connectionToCalendarAPIOn)
	s.Step(`^There is the event:$`, thereIsTheEvent)
	s.Step(`^I send AddEvent request$`, iSendAddEventRequest)
	s.Step(`^response should have event$`, responseShouldHaveEvent)

	s.AfterScenario(closeClient)
}
