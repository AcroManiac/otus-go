package storage

import (
	"testing"
	"time"

	"github.com/AcroManiac/otus-go/homework/calendar/internal/event"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

func TestMemoryStorage_Add(t *testing.T) {
	var addTests = struct {
		expected []time.Time
	}{
		[]time.Time{events[0].StartTime, events[1].StartTime, events[2].StartTime, events[3].StartTime},
	}
	ms := createMemoryStorage(t)
	for _, startTime := range addTests.expected {
		if !ms.isExist(startTime) {
			t.Errorf("event with time %v was not populated", startTime)
		}
	}

	// Test error case
	err := ms.Add(events[0])
	assert.NotNil(t, err, "Method should return an error")
}

func TestMemoryStorage_Remove(t *testing.T) {
	var removeTests = struct {
		existed []time.Time
		removed []time.Time
	}{
		[]time.Time{events[0].StartTime, events[2].StartTime},
		[]time.Time{events[1].StartTime, events[3].StartTime},
	}
	ms := createMemoryStorage(t)

	// Remove events from storage
	for _, startTime := range removeTests.removed {
		err := ms.Remove(startTime)
		assert.Nil(t, err, "Method should return no error")

		// Check if event was removed really
		ok := ms.isExist(startTime)
		assert.Falsef(t, ok, "event with time %v shouldn't exist in storage", startTime)
	}

	// Check remained events
	for _, startTime := range removeTests.existed {
		ok := ms.isExist(startTime)
		assert.Truef(t, ok, "event with time %v should exist in storage", startTime)
	}
}

func TestMemoryStorage_Edit(t *testing.T) {
	ms := createMemoryStorage(t)

	// Modify time for event in storage
	me := events[3]
	me.StartTime = events[0].StartTime.Add(time.Hour)
	err := ms.Edit(events[3].StartTime, me)
	assert.Nil(t, err, "Method should return no error")
	me.StartTime = events[0].StartTime
	err = ms.Edit(events[3].StartTime, me)
	assert.NotNil(t, err, "Method should return error")
}

func TestMemoryStorage_GetByTimePeriod(t *testing.T) {
	var getTests = struct {
		expected []event.Event
	}{
		events[0:3],
	}
	ms := createMemoryStorage(t)

	filtered, err := ms.GetByTimePeriod(events[0].StartTime, Week)
	assert.Nil(t, err, "Method should return no error")
	assert.ElementsMatch(t, getTests.expected, filtered, "Events should match")
}
