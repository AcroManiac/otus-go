package interfaces

import "github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/entities"

// EventsCollector interface
type EventsCollector interface {
	// Function gets events from storage.
	// Returns slice of events or error if failed
	GetEvents() ([]entities.Event, error)
}
