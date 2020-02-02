package calendar

import (
	"time"

	"github.com/AcroManiac/otus-go/homework/calendar/internal/event"
	"github.com/AcroManiac/otus-go/homework/calendar/internal/storage"
)

type Calendar interface {
	CreateEvent(startTime time.Time, stopTime time.Time) (storage.EventId, error)
	EditEvent(id storage.EventId, event event.Event) error
	DeleteEvent(id storage.EventId) error
}
