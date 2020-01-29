package storage

import (
	"github.com/AcroManiac/otus-go/homework/calendar/internal/event"
	"time"
)

type MemoryStorage struct {
	events map[time.Time]event.Event
}

func New() *MemoryStorage {
	return &MemoryStorage{events: make(map[time.Time]event.Event)}
}

func (ms *MemoryStorage) Add(event event.Event) error {
	if _, ok := ms.events[event.StartTime]; ok {
		return ErrTimeBusy
	}
	ms.events[event.StartTime] = event
	return nil
}

// isExist is checking element existence in map.
// Used for unit testing purposes also
func (ms MemoryStorage) isExist(time time.Time) bool {
	_, ok := ms.events[time]
	return ok
}

func (ms *MemoryStorage) Remove(time time.Time) error {
	if !ms.isExist(time) {
		return ErrNotFoundEvent
	}
	delete(ms.events, time)
	return nil
}

func (ms *MemoryStorage) Edit(time time.Time, event event.Event) error {
	// Check input data for errors
	if !ms.isExist(time) {
		return ErrNotFoundEvent
	}

	// Check if event time was changed
	if time != event.StartTime {
		if ms.isExist(event.StartTime) {
			return ErrTimeBusy
		}
		ms.Remove(time)
	}
	ms.events[time] = event
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
	for key, value := range ms.events {
		if key.After(startTime) && key.Before(stopTime) {
			selected = append(selected, value)
		}
	}

	return selected, nil
}
