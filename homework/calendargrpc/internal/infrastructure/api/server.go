package api

import (
	"context"
	"errors"
	"time"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/entities"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/logger"
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes"

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
	logger.Debug("Received CreateEvent request", "content", request.String())

	// Convert time from gRPC representation
	var startTime time.Time
	if request.GetStartTime() != nil {
		st, err := ptypes.Timestamp(request.GetStartTime())
		if err != nil {
			logger.Error("failed to convert timestamp", "error", err)
			return nil, err
		}
		startTime = st
	}

	var duration time.Duration
	if request.GetDuration() != nil {
		d, err := ptypes.Duration(request.GetDuration())
		if err != nil {
			logger.Error("failed to convert duration", "error", err)
			return nil, err
		}
		duration = d
	}

	// Create new calendar event
	id, err := c.cal.CreateEvent(
		request.GetTitle(),
		request.GetDescription(),
		request.GetOwner(),
		startTime,
		duration)
	if err != nil {
		// Send error response
		response := &CreateEventResponse{
			Result: &CreateEventResponse_Error{
				Error: err.Error(),
			},
		}
		logger.Error("error creating new calendar event", "error", err)
		return response, err
	}

	// Create output protobuf message
	message := &Event{
		Id:          uuid.UUID(id).String(),
		Title:       request.GetTitle(),
		Description: request.GetDescription(),
		Owner:       request.GetOwner(),
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
	logger.Debug("Sending CreateEvent response", "content", response.String())
	return response, nil
}

func (c *CalendarApiServerImpl) EditEvent(
	ctx context.Context, request *EditEventRequest) (*EditEventResponse, error) {
	logger.Debug("Received EditEvent request", "content", request.String())

	input := request.GetEvent()

	// Convert time from gRPC representation
	var startTime time.Time
	if input.GetStartTime() != nil {
		st, err := ptypes.Timestamp(input.GetStartTime())
		if err != nil {
			logger.Error("failed to convert timestamp", "error", err)
			return nil, err
		}
		startTime = st
	}

	var duration time.Duration
	if input.GetDuration() != nil {
		d, err := ptypes.Duration(input.GetDuration())
		if err != nil {
			logger.Error("failed to convert duration", "error", err)
			return nil, err
		}
		duration = d
	}

	var notify time.Duration
	if input.GetNotify() != nil {
		n, err := ptypes.Duration(input.GetNotify())
		if err != nil {
			logger.Error("failed to convert duration", "error", err)
			return nil, err
		}
		notify = n
	}

	idTemp, err := uuid.FromString(input.GetId())
	if err != nil {
		logger.Error("failed to convert uuid", "error", err)
		return nil, err
	}
	id := entities.IdType(idTemp)
	description := input.GetDescription()
	ev := entities.Event{
		Id:          id,
		Title:       input.GetTitle(),
		StartTime:   startTime,
		Duration:    duration,
		Description: &description,
		Owner:       input.GetOwner(),
		Notify:      &notify,
	}

	// Edit calendar event
	err = c.cal.EditEvent(id, ev)
	if err != nil {
		logger.Error("failed to edit calendar event", "error", err)
	}

	// Send output response
	response := &EditEventResponse{
		Error: func() string {
			if err != nil {
				return err.Error()
			}
			return ""
		}(),
	}
	logger.Debug("Sending EditEvent response", "content", response.String())
	return response, nil
}

func (c *CalendarApiServerImpl) DeleteEvent(
	ctx context.Context, request *DeleteEventRequest) (*DeleteEventResponse, error) {
	logger.Debug("Received DeleteEvent request", "content", request.String())

	// Get Id from request
	idTemp, err := uuid.FromString(request.GetId())
	if err != nil {
		logger.Error("failed to convert uuid", "error", err)
		return nil, err
	}
	id := entities.IdType(idTemp)

	// Delete calendar event
	err = c.cal.DeleteEvent(id)
	if err != nil {
		logger.Error("failed to delete calendar event", "error", err)
	}

	// Send output response
	response := &DeleteEventResponse{
		Error: func() string {
			if err != nil {
				return err.Error()
			}
			return ""
		}(),
	}
	logger.Debug("Sending DeleteEvent response", "content", response.String())
	return response, nil
}

func (c *CalendarApiServerImpl) GetEvents(
	ctx context.Context, request *GetEventsRequest) (*GetEventsResponse, error) {
	logger.Debug("Received GetEvents request", "content", request.String())

	// Convert time from gRPC representation
	var startTime time.Time
	if request.GetStartTime() != nil {
		st, err := ptypes.Timestamp(request.GetStartTime())
		if err != nil {
			logger.Error("failed to convert timestamp", "error", err)
			return nil, err
		}
		startTime = st
	}

	// Get and convert time period
	var period entities.TimePeriod
	switch request.GetPeriod() {
	case TimePeriod_TIME_DAY:
		period = entities.Day
	case TimePeriod_TIME_WEEK:
		period = entities.Week
	case TimePeriod_TIME_MONTH:
		period = entities.Month
	case TimePeriod_TIME_UNKNOWN:
		err := errors.New("wrong input time period")
		logger.Error(err.Error())
		return nil, err
	}

	// Get calendar events
	events, err := c.cal.GetEventsByTimePeriod(period, startTime)
	if err != nil {
		logger.Error("failed to get calendar events", "error", err)
	}

	// Fill output data
	var grpcEvents []*Event
	for _, ev := range events {
		// Convert time to gRPC representation
		startTime, err := ptypes.TimestampProto(ev.StartTime)
		if err != nil {
			logger.Error("failed to convert timestamp", "error", err)
			return nil, err
		}

		duration := ptypes.DurationProto(ev.Duration)
		notify := ptypes.DurationProto(time.Microsecond)
		if ev.Notify != nil {
			notify = ptypes.DurationProto(*ev.Notify)
		}

		grpcEvents = append(grpcEvents, &Event{
			Id:    uuid.UUID(ev.Id).String(),
			Title: ev.Title,
			Description: func(s *string) string {
				if s != nil {
					return *s
				}
				return ""
			}(ev.Description),
			Owner:     ev.Owner,
			StartTime: startTime,
			Duration:  duration,
			Notify:    notify,
		})
	}

	// Send output response
	response := &GetEventsResponse{
		Events: grpcEvents,
		Error: func() string {
			if err != nil {
				return err.Error()
			}
			return ""
		}(),
	}
	logger.Debug("Sending GetEvents response", "content", response.String())
	return response, nil
}
