package interfaces

import (
	"time"

	"github.com/ahamtat/otus-go/homework/calendargrpc/internal/domain/entities"
)

// Storage interfaces is used for event manipulation in
// a data storage (memory, database, etc...)
type Storage interface {
	// Add event to data storage.
	// The function returns event Id or error if event time is occupied
	Add(event entities.Event) (entities.IdType, error)

	// Remove event from data storage.
	// If there is no event with time specified the function returns error
	Remove(id entities.IdType) error

	// Edit event data in data storage
	// If event time is occupied or not in storage function returns error
	Edit(id entities.IdType, ev entities.Event) error

	// Get event slice by time period (Day, Week, Month)
	// Function estimates start and stop of specified time period
	// and returns events fit in it
	GetEventsByTimePeriod(period entities.TimePeriod, time time.Time) ([]entities.Event, error)
}
