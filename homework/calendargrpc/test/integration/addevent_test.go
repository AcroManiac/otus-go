package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/entities"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/pkg/api"
	"github.com/cucumber/messages-go/v10"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
)

// Data for event addition
type addEventTest struct {
	eventData      *api.CreateEventRequest
	createResponse *api.CreateEventResponse
}

func (t *addEventTest) thereIsTheEvent(arg1 *messages.PickleStepArgument_PickleDocString) error {
	if arg1 == nil {
		return errors.New("argument is invalid")
	}

	event := &entities.Event{}
	if err := json.Unmarshal([]byte(arg1.Content), event); err != nil {
		return errors.Wrap(err, "couldn't parse JSON object")
	}

	startTime, err := ptypes.TimestampProto(event.StartTime)
	if err != nil {
		return errors.Wrap(err, "error converting timestamp")
	}

	t.eventData = &api.CreateEventRequest{
		Title:       event.Title,
		Description: *event.Description,
		Owner:       event.Owner,
		StartTime:   startTime,
		Duration:    ptypes.DurationProto(time.Minute),
		Notify:      ptypes.DurationProto(time.Second),
	}
	return nil
}

func (t *addEventTest) iSendAddEventRequest() error {

	// Create cancel context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Send create event request to gRPC server
	var err error
	t.createResponse, err = grpcClient.CreateEvent(ctx, t.eventData)
	if err != nil {
		return errors.Wrap(err, "failed sending CreateEvent request")
	}

	return nil
}

func (t *addEventTest) responseShouldHaveEvent() error {

	// Get new event id from gRPC response
	respEvent := t.createResponse.GetEvent()
	if respEvent == nil {
		return errors.New("response returned no event")
	}
	if len(respEvent.GetId()) == 0 {
		return errors.New("no valid event ID returned")
	}
	return nil
}
