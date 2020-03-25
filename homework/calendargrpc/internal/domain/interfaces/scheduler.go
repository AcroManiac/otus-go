package interfaces

import (
	"time"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/entities"
)

// Scheduler interface
type Scheduler interface {
	// Schedule events to be processed by services
	Schedule() error

	// Clean retained events
	Clean() error
}

// Cleaner interface for objects cleaning
type Cleaner interface {
	// Clean objects
	Clean() error
}

// RetentionPolicy interface for keeping objects in storage
type RetentionPolicy interface {
	// GetDuration returns retention period duration
	GetDuration() time.Duration
}

// EventsCollector interface
type EventsCollector interface {
	// Function gets events from storage.
	// Returns slice of events or error if failed
	GetEvents() ([]entities.Event, error)
}
