package calendar

import (
	"log"
	"time"

	"github.com/AcroManiac/otus-go/homework/calendar/internal/event"
	"github.com/AcroManiac/otus-go/homework/calendar/internal/storage"
)

type Impl struct {
	storage storage.Storage
}

func NewCalendar() *Impl {
	return &Impl{storage: storage.NewMemoryStorage()}
}

func (i *Impl) CreateEvent(startTime time.Time, stopTime time.Time) (storage.EventId, error) {
	// Create and add event
	newEvent := event.Event{
		Header:       "",
		StartTime:    startTime,
		StopTime:     stopTime,
		Description:  nil,
		Owner:        "",
		Notification: nil,
	}

	id, err := i.storage.Add(newEvent)
	if err != nil {
		log.Printf("Error adding event to storage: %s", err.Error())
		return id, err
	}

	return id, nil
}

func (i *Impl) EditEvent(id storage.EventId, event event.Event) error {
	return i.storage.Edit(id, event)
}

func (i *Impl) DeleteEvent(id storage.EventId) error {
	return i.storage.Remove(id)
}
