package calendar

import (
	"time"

	"github.com/AcroManiac/otus-go/homework/calendar/internal/event"
	"github.com/AcroManiac/otus-go/homework/calendar/internal/storage"
)

// Calendar interface
type Calendar interface {

	// CreateEvent function constructs new event with default values
	// and adds it to event storage. Returns Id of generated event or
	// storage error if any
	CreateEvent(startTime time.Time, stopTime time.Time) (storage.EventId, error)

	// EditEvent finds event in storage by Id and replaces it
	// with event in parameters. Returns error if event Id is wrong
	EditEvent(id storage.EventId, event event.Event) error

	// DeleteEvent removes event from storage by Id.
	// Returns error if event Id is wrong
	DeleteEvent(id storage.EventId) error
}
