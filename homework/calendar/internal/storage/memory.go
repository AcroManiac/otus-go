package storage

import (
	"sync"
	"time"

	"github.com/ahamtat/otus-go/homework/calendar/internal/event"
	"github.com/google/uuid"
)

type MemoryStorage struct {
	mu     sync.RWMutex
	events map[event.IdType]event.Event
}

func NewStorage() Storage {
	return &MemoryStorage{events: make(map[event.IdType]event.Event)}
}

// isExistTime is checking event time existence in map.
// Used for unit testing purposes also
func (ms MemoryStorage) isExistTime(time time.Time) (event.IdType, bool) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	for id, e := range ms.events {
		if e.StartTime == time {
			return id, true
		}
	}
	return event.IdType(uuid.UUID{}), false
}

// isExistId is checking event existence in map by Id.
func (ms MemoryStorage) isExistId(id event.IdType) bool {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	_, ok := ms.events[id]
	return ok
}

func (ms *MemoryStorage) Add(ev event.Event) (event.IdType, error) {
	if _, ok := ms.isExistTime(ev.StartTime); ok {
		return event.IdType(uuid.UUID{}), ErrTimeBusy
	}
	id := event.IdType(uuid.New())

	ms.mu.Lock()
	ms.events[id] = ev
	ms.mu.Unlock()

	return id, nil
}

func (ms *MemoryStorage) Remove(id event.IdType) error {
	if !ms.isExistId(id) {
		return ErrNotFoundEvent
	}

	ms.mu.Lock()
	delete(ms.events, id)
	ms.mu.Unlock()

	return nil
}

func (ms *MemoryStorage) Edit(id event.IdType, ev event.Event) error {
	// Check input data for errors
	if !ms.isExistId(id) {
		return ErrNotFoundEvent
	}
	if _, ok := ms.isExistTime(ev.StartTime); ok {
		return ErrTimeBusy
	}

	ms.mu.Lock()
	ms.events[id] = ev
	ms.mu.Unlock()

	return nil
}

func (ms MemoryStorage) GetByTimePeriod(t time.Time, period event.TimePeriod) ([]event.Event, error) {
	var selected []event.Event
	var startTime, stopTime time.Time

	// Calculate start and stop times for searching interval
	switch period {
	case event.Day:
		startTime = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		stopTime = startTime.Add(24 * time.Hour)
	case event.Month:
		startTime = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
		stopTime = startTime.AddDate(0, 1, 0)
	case event.Week:
		startTime = t.AddDate(0, 0, -int(t.Weekday()))
		stopTime = startTime.AddDate(0, 0, 7)
	}

	// Iterate through map to find matching events
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	for _, e := range ms.events {
		if e.StartTime.After(startTime) && e.StartTime.Before(stopTime) {
			selected = append(selected, e)
		}
	}

	return selected, nil
}
