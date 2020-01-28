package storage

import (
	"github.com/AcroManiac/otus-go/homework/calendar/internal/event"
	"github.com/google/uuid"
	"testing"
	"time"
)

var location, _ = time.LoadLocation("Europe/Moscow")
var timePoint = time.Date(2020, 1, 28, 16, 20, 39, 50, location)
var events = []event.Event{
	{
		Id:           uuid.New(),
		Header:       "Test Event 1",
		StartTime:    timePoint,
		Duration:     time.Hour,
		Description:  func(s string) *string { return &s }("Memory storage test event 1"),
		Owner:        "artem",
		Notification: func(t time.Duration) *time.Duration { return &t }(15 * time.Minute),
	},
	{
		Id:           uuid.New(),
		Header:       "Test Event 2",
		StartTime:    timePoint.Add(6 * time.Hour),
		Duration:     30 * time.Minute,
		Description:  func(s string) *string { return &s }("Memory storage test event 2"),
		Owner:        "artem",
		Notification: nil,
	},
	{
		Id:           uuid.New(),
		Header:       "Test Event 3",
		StartTime:    timePoint.Add(-time.Hour),
		Duration:     15 * time.Minute,
		Description:  func(s string) *string { return &s }("Memory storage test event 3"),
		Owner:        "artem",
		Notification: func(t time.Duration) *time.Duration { return &t }(5 * time.Minute),
	},
	{
		Id:           uuid.New(),
		Header:       "Test Event 4",
		StartTime:    timePoint.Add(7 * 24 * time.Hour), // Next week
		Duration:     15 * time.Minute,
		Description:  func(s string) *string { return &s }("Memory storage test event 4"),
		Owner:        "artem",
		Notification: func(t time.Duration) *time.Duration { return &t }(5 * time.Minute),
	},
}

// Factory for memory storage. Build and populate with events
func createMemoryStorage(t *testing.T) *MemoryStorage {
	ms := New()
	for _, e := range events {
		if err := ms.Add(e); err != nil {
			t.Errorf("Couldn't populate with event: %s", err.Error())
		}
	}
	return ms
}

var addTests = struct {
	expected []time.Time
}{
	[]time.Time{events[0].StartTime, events[1].StartTime, events[2].StartTime, events[3].StartTime},
}

func TestMemoryStorage_Add(t *testing.T) {
	ms := createMemoryStorage(t)
	for _, startTime := range addTests.expected {
		if !ms.isExist(startTime) {
			t.Errorf("event with time %v was not populated", startTime)
		}
	}
}
