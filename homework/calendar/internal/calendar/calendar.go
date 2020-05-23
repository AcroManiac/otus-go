package calendar

import (
	"time"

	"github.com/ahamtat/otus-go/homework/calendar/internal/event"
)

// Calendar interface
type Calendar interface {

	// CreateEvent function constructs new event with default values
	// and adds it to event storage. Returns Id of generated event or
	// storage error if any
	CreateEvent(startTime time.Time, stopTime time.Time) (event.IdType, error)

	// EditEvent finds event in storage by Id and replaces it
	// with event in parameters. Returns error if event Id is wrong
	EditEvent(id event.IdType, ev event.Event) error

	// DeleteEvent removes event from storage by Id.
	// Returns error if event Id is wrong
	DeleteEvent(id event.IdType) error

	// Get event slice by time period (Day, Week, Month)
	// Returns empty slice if there are no events
	GetEventsByTimePeriod(time time.Time, period event.TimePeriod) ([]event.Event, error)
}
