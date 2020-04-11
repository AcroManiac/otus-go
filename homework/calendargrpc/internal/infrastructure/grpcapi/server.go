package grpcapi

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/monitoring"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/pkg/api"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/entities"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/logger"
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/interfaces"
)

type CalendarApiServerImpl struct {
	cal   interfaces.Calendar
	stats map[string]*prometheus.SummaryVec
	errs  *prometheus.CounterVec
}

func NewCalendarApiServer(cal interfaces.Calendar) api.CalendarApiServer {
	return &CalendarApiServerImpl{
		cal: cal,
		stats: map[string]*prometheus.SummaryVec{
			"create": monitoring.NewSummaryVec("calendar_api", "CreateEvent", "Create event statistics"),
			"edit":   monitoring.NewSummaryVec("calendar_api", "EditEvent", "Edit event statistics"),
			"delete": monitoring.NewSummaryVec("calendar_api", "DeleteEvent", "Delete event statistics"),
			"get":    monitoring.NewSummaryVec("calendar_api", "GetEvents", "Get events statistics"),
		},
		errs: monitoring.NewErrorVec("calendar_api"),
	}
}

func (c *CalendarApiServerImpl) CreateEvent(
	ctx context.Context, request *api.CreateEventRequest) (*api.CreateEventResponse, error) {
	logger.Debug("Received CreateEvent request", "content", request.String())

	// Start function execution time counting
	startFunc := time.Now()

	// Convert time from gRPC representation
	var startTime time.Time
	if request.GetStartTime() != nil {
		st, err := ptypes.Timestamp(request.GetStartTime())
		if err != nil {
			logger.Error("failed to convert timestamp", "error", err)
			c.errs.WithLabelValues(codes.InvalidArgument.String())
			return nil, status.Error(codes.InvalidArgument, "failed to convert timestamp")
		}
		startTime = st
	}

	var duration time.Duration
	if request.GetDuration() != nil {
		d, err := ptypes.Duration(request.GetDuration())
		if err != nil {
			logger.Error("failed to convert timestamp", "error", err)
			c.errs.WithLabelValues(codes.InvalidArgument.String())
			return nil, status.Error(codes.InvalidArgument, "failed to convert timestamp")
		}
		duration = d
	}

	var notify time.Duration
	if request.GetDuration() != nil {
		n, err := ptypes.Duration(request.GetNotify())
		if err != nil {
			logger.Error("failed to convert duration", "error", err)
			c.errs.WithLabelValues(codes.InvalidArgument.String())
			return nil, status.Error(codes.InvalidArgument, "failed to convert duration")
		}
		notify = n
	}

	// Create new calendar event
	id, err := c.cal.CreateEvent(
		request.GetTitle(),
		request.GetDescription(),
		request.GetOwner(),
		startTime,
		duration,
		notify)
	if err != nil {
		// Send error response
		response := &api.CreateEventResponse{
			Result: &api.CreateEventResponse_Error{
				Error: err.Error(),
			},
		}
		logger.Error("error creating new calendar event", "error", err)
		c.errs.WithLabelValues(codes.AlreadyExists.String())
		return response, status.Error(codes.AlreadyExists, "error creating new calendar event")
	}

	// Create output protobuf message
	message := &api.Event{
		Id:          uuid.UUID(id).String(),
		Title:       request.GetTitle(),
		Description: request.GetDescription(),
		Owner:       request.GetOwner(),
		StartTime:   request.GetStartTime(),
		Duration:    request.GetDuration(),
		Notify:      request.GetNotify(),
	}

	// Send data response
	response := &api.CreateEventResponse{
		Result: &api.CreateEventResponse_Event{
			Event: message,
		},
	}

	// Store duration of function execution
	dur := time.Since(startFunc)
	c.stats["create"].WithLabelValues("duration").Observe(dur.Seconds())

	logger.Debug("Sending CreateEvent response", "content", response.String())
	c.errs.WithLabelValues(codes.OK.String())
	return response, nil
}

func (c *CalendarApiServerImpl) EditEvent(
	ctx context.Context, request *api.EditEventRequest) (*api.EditEventResponse, error) {
	logger.Debug("Received EditEvent request", "content", request.String())

	// Start function execution time counting
	startFunc := time.Now()

	input := request.GetEvent()

	// Convert time from gRPC representation
	var startTime time.Time
	if input.GetStartTime() != nil {
		st, err := ptypes.Timestamp(input.GetStartTime())
		if err != nil {
			logger.Error("failed to convert timestamp", "error", err)
			c.errs.WithLabelValues(codes.InvalidArgument.String())
			return nil, status.Error(codes.InvalidArgument, "failed to convert timestamp")
		}
		startTime = st
	}

	var duration time.Duration
	if input.GetDuration() != nil {
		d, err := ptypes.Duration(input.GetDuration())
		if err != nil {
			logger.Error("failed to convert timestamp", "error", err)
			c.errs.WithLabelValues(codes.InvalidArgument.String())
			return nil, status.Error(codes.InvalidArgument, "failed to convert timestamp")
		}
		duration = d
	}

	var notify time.Duration
	if input.GetNotify() != nil {
		n, err := ptypes.Duration(input.GetNotify())
		if err != nil {
			logger.Error("failed to convert duration", "error", err)
			c.errs.WithLabelValues(codes.InvalidArgument.String())
			return nil, status.Error(codes.InvalidArgument, "failed to convert duration")
		}
		notify = n
	}

	idTemp, err := uuid.FromString(input.GetId())
	if err != nil {
		logger.Error("failed to convert uuid", "error", err)
		c.errs.WithLabelValues(codes.InvalidArgument.String())
		return nil, status.Error(codes.InvalidArgument, "failed to convert uuid")
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
	response := &api.EditEventResponse{
		Error: func() string {
			if err != nil {
				return err.Error()
			}
			return ""
		}(),
	}

	// Store duration of function execution
	dur := time.Since(startFunc)
	c.stats["edit"].WithLabelValues("duration").Observe(dur.Seconds())

	logger.Debug("Sending EditEvent response", "content", response.String())
	c.errs.WithLabelValues(codes.OK.String())
	return response, nil
}

func (c *CalendarApiServerImpl) DeleteEvent(
	ctx context.Context, request *api.DeleteEventRequest) (*api.DeleteEventResponse, error) {
	logger.Debug("Received DeleteEvent request", "content", request.String())

	// Start function execution time counting
	startFunc := time.Now()

	// Get Id from request
	idTemp, err := uuid.FromString(request.GetId())
	if err != nil {
		logger.Error("failed to convert uuid", "error", err)
		c.errs.WithLabelValues(codes.InvalidArgument.String())
		return nil, status.Error(codes.InvalidArgument, "failed to convert uuid")
	}
	id := entities.IdType(idTemp)

	// Delete calendar event
	err = c.cal.DeleteEvent(id)
	if err != nil {
		logger.Error("failed to delete calendar event", "error", err)
	}

	// Send output response
	response := &api.DeleteEventResponse{
		Error: func() string {
			if err != nil {
				return err.Error()
			}
			return ""
		}(),
	}

	// Store duration of function execution
	dur := time.Since(startFunc)
	c.stats["delete"].WithLabelValues("duration").Observe(dur.Seconds())

	logger.Debug("Sending DeleteEvent response", "content", response.String())
	c.errs.WithLabelValues(codes.OK.String())
	return response, nil
}

func (c *CalendarApiServerImpl) GetEvents(
	ctx context.Context, request *api.GetEventsRequest) (*api.GetEventsResponse, error) {
	logger.Debug("Received GetEvents request", "content", request.String())

	// Start function execution time counting
	startFunc := time.Now()

	// Convert time from gRPC representation
	var startTime time.Time
	if request.GetStartTime() != nil {
		st, err := ptypes.Timestamp(request.GetStartTime())
		if err != nil {
			logger.Error("failed to convert timestamp", "error", err)
			c.errs.WithLabelValues(codes.InvalidArgument.String())
			return nil, status.Error(codes.InvalidArgument, "failed to convert timestamp")
		}
		startTime = st
	}

	// Get and convert time period
	var period entities.TimePeriod
	switch request.GetPeriod() {
	case api.TimePeriod_TIME_DAY:
		period = entities.Day
	case api.TimePeriod_TIME_WEEK:
		period = entities.Week
	case api.TimePeriod_TIME_MONTH:
		period = entities.Month
	case api.TimePeriod_TIME_UNKNOWN:
		err := errors.New("wrong input time period")
		logger.Error(err.Error())
		c.errs.WithLabelValues(codes.InvalidArgument.String())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Get calendar events
	events, err := c.cal.GetEventsByTimePeriod(period, startTime)
	if err != nil {
		logger.Error("failed to get calendar events", "error", err)
	}

	// Fill output data
	var grpcEvents []*api.Event
	for _, ev := range events {
		// Convert time to gRPC representation
		startTime, err := ptypes.TimestampProto(ev.StartTime)
		if err != nil {
			logger.Error("failed to convert timestamp", "error", err)
			c.errs.WithLabelValues(codes.InvalidArgument.String())
			return nil, status.Error(codes.InvalidArgument, "failed to convert timestamp")
		}

		duration := ptypes.DurationProto(ev.Duration)
		notify := ptypes.DurationProto(time.Microsecond)
		if ev.Notify != nil {
			notify = ptypes.DurationProto(*ev.Notify)
		}

		grpcEvents = append(grpcEvents, &api.Event{
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
	response := &api.GetEventsResponse{
		Events: grpcEvents,
		Error: func() string {
			if err != nil {
				return err.Error()
			}
			return ""
		}(),
	}

	// Store duration of function execution
	dur := time.Since(startFunc)
	c.stats["get"].WithLabelValues("duration").Observe(dur.Seconds())

	logger.Debug("Sending GetEvents response", "content", response.String())
	c.errs.WithLabelValues(codes.OK.String())
	return response, nil
}
