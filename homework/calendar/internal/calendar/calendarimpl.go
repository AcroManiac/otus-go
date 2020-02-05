package calendar

import (
	"time"

	"github.com/AcroManiac/otus-go/homework/calendar/internal/logger"

	"github.com/AcroManiac/otus-go/homework/calendar/internal/event"
	"github.com/AcroManiac/otus-go/homework/calendar/internal/storage"
)

type Impl struct {
	storage storage.Storage
}

func NewCalendar() *Impl {
	return &Impl{storage: storage.NewMemoryStorage()}
}

func (i *Impl) CreateEvent(startTime time.Time, stopTime time.Time) (event.IdType, error) {
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
		logger.Error("Error adding event to storage", "error", err)
		return id, err
	}

	return id, nil
}

func (i *Impl) EditEvent(id event.IdType, ev event.Event) error {
	return i.storage.Edit(id, ev)
}

func (i *Impl) DeleteEvent(id event.IdType) error {
	return i.storage.Remove(id)
}

func (i *Impl) GetEventsByTimePeriod(time time.Time, period event.TimePeriod) ([]event.Event, error) {
	return i.storage.GetByTimePeriod(time, period)
}
