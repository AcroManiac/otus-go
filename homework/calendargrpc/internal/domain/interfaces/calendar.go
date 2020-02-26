package interfaces

import (
	"time"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/entities"
)

// Calendar interfaces
type Calendar interface {

	// CreateEvent function constructs new event with default values
	// and adds it to event storage. Returns Id of generated event or
	// storage error if any
	CreateEvent(title, description string, startTime time.Time, duration time.Duration) (entities.IdType, error)

	// EditEvent finds event in storage by Id and replaces it
	// with event in parameters. Returns error if event Id is wrong
	EditEvent(id entities.IdType, ev entities.Event) error

	// DeleteEvent removes event from storage by Id.
	// Returns error if event Id is wrong
	DeleteEvent(id entities.IdType) error

	// Get event slice by time period (Day, Week, Month)
	// Returns empty slice if there are no events
	GetEventsByTimePeriod(period entities.TimePeriod, time time.Time) ([]entities.Event, error)
}
