package storage

import (
	"time"

	"github.com/AcroManiac/otus-go/homework/calendar/internal/event"
	"github.com/google/uuid"
)

type MemoryStorage struct {
	events map[EventId]event.Event
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{events: make(map[EventId]event.Event)}
}

// isExistTime is checking event time existence in map.
// Used for unit testing purposes also
func (ms MemoryStorage) isExistTime(time time.Time) (EventId, bool) {
	for id, e := range ms.events {
		if e.StartTime == time {
			return id, true
		}
	}
	return EventId(uuid.UUID{}), false
}

// isExistId is checking event existence in map by Id.
func (ms MemoryStorage) isExistId(id EventId) bool {
	_, ok := ms.events[id]
	return ok
}

func (ms *MemoryStorage) Add(event event.Event) (EventId, error) {
	if _, ok := ms.isExistTime(event.StartTime); ok {
		return EventId(uuid.UUID{}), ErrTimeBusy
	}
	eventId := EventId(uuid.New())
	ms.events[eventId] = event
	return eventId, nil
}

func (ms *MemoryStorage) Remove(id EventId) error {
	if !ms.isExistId(id) {
		return ErrNotFoundEvent
	}
	delete(ms.events, id)
	return nil
}

func (ms *MemoryStorage) Edit(id EventId, event event.Event) error {
	// Check input data for errors
	if !ms.isExistId(id) {
		return ErrNotFoundEvent
	}
	if _, ok := ms.isExistTime(event.StartTime); ok {
		return ErrTimeBusy
	}
	ms.events[id] = event
	return nil
}

func (ms MemoryStorage) GetByTimePeriod(t time.Time, period TimePeriod) ([]event.Event, error) {
	var selected []event.Event
	var startTime, stopTime time.Time

	// Calculate start and stop times for searching interval
	switch period {
	case Day:
		startTime = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		stopTime = startTime.Add(24 * time.Hour)
	case Month:
		startTime = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
		stopTime = startTime.AddDate(0, 1, 0)
	case Week:
		startTime = t.AddDate(0, 0, -int(t.Weekday()))
		stopTime = startTime.AddDate(0, 0, 7)
	}

	// Iterate through map to find matching events
	for _, e := range ms.events {
		if e.StartTime.After(startTime) && e.StartTime.Before(stopTime) {
			selected = append(selected, e)
		}
	}

	return selected, nil
}
