package storage

import (
	"time"

	"github.com/AcroManiac/otus-go/homework/calendar/internal/event"
)

// Storage interface is used for event manipulation in
// a data storage (memory, database, etc...)
type Storage interface {
	// Add event to data storage.
	// The function returns event Id or error if event time is occupied
	Add(event event.Event) (event.IdType, error)

	// Remove event from data storage.
	// If there is no event with time specified the function returns error
	Remove(id event.IdType) error

	// Edit event data in data storage
	// If event time is occupied or not in storage function returns error
	Edit(id event.IdType, ev event.Event) error

	// Get event slice by time period (Day, Week, Month)
	// Function estimates start and stop of specified time period
	// and returns events fit in it
	GetByTimePeriod(time time.Time, period event.TimePeriod) ([]event.Event, error)
}
