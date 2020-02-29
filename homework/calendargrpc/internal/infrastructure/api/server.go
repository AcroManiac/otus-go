package api

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"time"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/interfaces"
)

type CalendarApiServerImpl struct {
	cal interfaces.Calendar
}

func NewCalendarApiServer(cal interfaces.Calendar) CalendarApiServer {
	return &CalendarApiServerImpl{cal: cal}
}

func (c *CalendarApiServerImpl) CreateEvent(
	ctx context.Context, request *CreateEventRequest) (*CreateEventResponse, error) {
	// Convert time from gRPC representation
	var startTime time.Time
	if request.GetStartTime() != nil {
		st, err := ptypes.Timestamp(request.GetStartTime())
		if err != nil {
			return nil, err
		}
		startTime = &st
	}

	var duration time.Duration
	if request.GetDuration() != nil {
		d, err := ptypes.Duration(request.GetDuration())
		if err != nil {
			return nil, err
		}
		duration = &d
	}

	// Create new calendar event
	id, err := c.cal.CreateEvent(
		request.GetTitle(),
		request.GetDescription(),
		startTime,
		duration)
	if err != nil {
		// Send error response
		response := &CreateEventResponse{
			Result: &CreateEventResponse_Error{
				Error: err.Error(),
			},
		}
		return response, err
	}

	// Create output protobuf message
	message := &Event{
		Id:          id,
		Title:       request.GetTitle(),
		Description: request.GetDescription(),
		Owner:       "",
		StartTime:   request.GetStartTime(),
		Duration:    request.GetDuration(),
		Notify:      nil,
	}

	// Send data response
	response := &CreateEventResponse{
		Result: &CreateEventResponse_Event{
			Event: message,
		},
	}
	return response, nil
}

func (c *CalendarApiServerImpl) EditEvent(
	ctx context.Context, request *EditEventRequest) (*EditEventResponse, error) {
	//
	return nil, nil
}

func (c *CalendarApiServerImpl) DeleteEvent(
	ctx context.Context, request *DeleteEventRequest) (*DeleteEventResponse, error) {
	//
	return nil, nil
}

func (c *CalendarApiServerImpl) GetEvents(
	ctx context.Context, request *GetEventsRequest) (*GetEventsResponse, error) {
	//
	return nil, nil
}
