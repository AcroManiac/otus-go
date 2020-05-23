package logic

import (
	"time"

	"github.com/ahamtat/otus-go/homework/calendargrpc/internal/domain/entities"
	"github.com/ahamtat/otus-go/homework/calendargrpc/internal/domain/interfaces"
	"github.com/ahamtat/otus-go/homework/calendargrpc/internal/infrastructure/logger"
)

type Impl struct {
	storage interfaces.Storage
}

func NewCalendar(storage interfaces.Storage) interfaces.Calendar {
	return &Impl{storage: storage}
}

func (i *Impl) CreateEvent(
	title, description, owner string,
	startTime time.Time, duration time.Duration, notify time.Duration) (entities.IdType, error) {
	// Create and add event
	newEvent := entities.Event{
		Title:       title,
		StartTime:   startTime,
		Duration:    duration,
		Description: &description,
		Owner:       owner,
		Notify:      &notify,
	}

	id, err := i.storage.Add(newEvent)
	if err != nil {
		logger.Error("error adding event to storage", "error", err)
		return id, err
	}

	return id, nil
}

func (i *Impl) EditEvent(id entities.IdType, ev entities.Event) error {
	return i.storage.Edit(id, ev)
}

func (i *Impl) DeleteEvent(id entities.IdType) error {
	return i.storage.Remove(id)
}

func (i *Impl) GetEventsByTimePeriod(period entities.TimePeriod, time time.Time) ([]entities.Event, error) {
	return i.storage.GetEventsByTimePeriod(period, time)
}
