package logic

import (
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/entities"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/interfaces"
	"github.com/stretchr/testify/assert"
)

var location, _ = time.LoadLocation("Europe/Moscow")
var timePoint = time.Date(2020, 2, 2, 12, 12, 25, 0, location)
var events = []entities.Event{
	{
		Title:       "Test Event 1",
		StartTime:   timePoint,
		Duration:    time.Hour,
		Description: func(s string) *string { return &s }("Memory storage test event 1"),
		Owner:       "artem",
		Notify:      func(t time.Duration) *time.Duration { return &t }(15 * time.Minute),
	},
	{
		Title:       "Test Event 2",
		StartTime:   timePoint.Add(6 * time.Hour),
		Duration:    30 * time.Minute,
		Description: func(s string) *string { return &s }("Memory storage test event 2"),
		Owner:       "artem",
		Notify:      nil,
	},
	{
		Title:       "Test Event 3",
		StartTime:   timePoint.Add(-time.Hour),
		Duration:    15 * time.Minute,
		Description: func(s string) *string { return &s }("Memory storage test event 3"),
		Owner:       "artem",
		Notify:      func(t time.Duration) *time.Duration { return &t }(5 * time.Minute),
	},
	{
		Title:       "Test Event 4",
		StartTime:   timePoint.Add(7 * 24 * time.Hour), // Next week
		Duration:    15 * time.Minute,
		Description: func(s string) *string { return &s }("Memory storage test event 4"),
		Owner:       "artem",
		Notify:      func(t time.Duration) *time.Duration { return &t }(5 * time.Minute),
	},
}

// Factory for calendar. Build and populate with events
func createCalendar(t *testing.T) interfaces.Calendar {
	cal := NewCalendar(NewMemoryStorage())
	for _, e := range events {
		if _, err := cal.CreateEvent(e.Title, *e.Description, e.Owner, e.StartTime, e.Duration); err != nil {
			t.Errorf("Couldn't populate with event: %s", err.Error())
		}
	}
	return cal
}

func TestImpl_CreateEvent(t *testing.T) {
	cal := createCalendar(t)
	assert.NotNil(t, cal, "Object should not be nil")
}

func TestImpl_EditEvent(t *testing.T) {
	cal := createCalendar(t)
	assert.NotNil(t, cal, "Object should not be nil")

	// Create new event and add it to storage
	id, err := cal.CreateEvent(
		events[0].Title, *events[0].Description, events[0].Owner,
		events[0].StartTime.Add(time.Hour), 30*time.Minute)
	assert.Nil(t, err, "Error should be nil")

	// Modify existing event
	e := events[0]
	e.StartTime = e.StartTime.Add(2 * time.Hour)
	err = cal.EditEvent(id, e)
	assert.Nil(t, err, "Error should be nil")
}

func TestImpl_DeleteEvent(t *testing.T) {
	cal := createCalendar(t)
	assert.NotNil(t, cal, "Object should not be nil")

	// Create new event and add it to storage
	id, err := cal.CreateEvent(
		events[0].Title, *events[0].Description, events[0].Owner,
		events[0].StartTime.Add(time.Hour), 30*time.Minute)
	assert.Nil(t, err, "Error should be nil")

	// Delete existing event
	err = cal.DeleteEvent(id)
	assert.Nil(t, err, "Error should be nil")
}

func TestImpl_GetEventsByTimePeriod(t *testing.T) {
	cal := createCalendar(t)
	assert.NotNil(t, cal, "Object should not be nil")

	filtered, err := cal.GetEventsByTimePeriod(entities.Week, events[0].StartTime)
	assert.Nil(t, err, "Method should return no error")
	assert.NotEmpty(t, filtered, "Returned slice should not be empty")
}

///////////////////////////////////////////////////////////////////////////////

type MemoryStorage struct {
	mu     sync.RWMutex
	events map[entities.IdType]entities.Event
}

func NewMemoryStorage() interfaces.Storage {
	return &MemoryStorage{events: make(map[entities.IdType]entities.Event)}
}

// isExistTime is checking event time existence in map.
// Used for unit testing purposes also
func (ms MemoryStorage) isExistTime(time time.Time) (entities.IdType, bool) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	for id, e := range ms.events {
		if e.StartTime == time {
			return id, true
		}
	}
	return entities.IdType(uuid.UUID{}), false
}

// isExistId is checking event existence in map by Id.
func (ms MemoryStorage) isExistId(id entities.IdType) bool {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	_, ok := ms.events[id]
	return ok
}

func (ms *MemoryStorage) Add(ev entities.Event) (entities.IdType, error) {
	if _, ok := ms.isExistTime(ev.StartTime); ok {
		return entities.IdType(uuid.UUID{}), entities.ErrTimeBusy
	}
	id := entities.IdType(uuid.New())

	ms.mu.Lock()
	ms.events[id] = ev
	ms.mu.Unlock()

	return id, nil
}

func (ms *MemoryStorage) Remove(id entities.IdType) error {
	if !ms.isExistId(id) {
		return entities.ErrNotFoundEvent
	}

	ms.mu.Lock()
	delete(ms.events, id)
	ms.mu.Unlock()

	return nil
}

func (ms *MemoryStorage) Edit(id entities.IdType, ev entities.Event) error {
	// Check input data for errors
	if !ms.isExistId(id) {
		return entities.ErrNotFoundEvent
	}
	if _, ok := ms.isExistTime(ev.StartTime); ok {
		return entities.ErrTimeBusy
	}

	ms.mu.Lock()
	ms.events[id] = ev
	ms.mu.Unlock()

	return nil
}

func (ms MemoryStorage) GetEventsByTimePeriod(period entities.TimePeriod, t time.Time) ([]entities.Event, error) {
	var selected []entities.Event
	var startTime, stopTime time.Time

	// Calculate start and stop times for searching interval
	switch period {
	case entities.Day:
		startTime = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		stopTime = startTime.Add(24 * time.Hour)
	case entities.Month:
		startTime = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
		stopTime = startTime.AddDate(0, 1, 0)
	case entities.Week:
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
