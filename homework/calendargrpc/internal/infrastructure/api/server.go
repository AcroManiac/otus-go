package api

import (
	"context"

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
	//
	return nil, nil
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
