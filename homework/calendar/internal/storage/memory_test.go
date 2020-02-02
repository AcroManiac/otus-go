package storage

import (
	"testing"
	"time"

	"github.com/AcroManiac/otus-go/homework/calendar/internal/event"
	"github.com/stretchr/testify/assert"
)

var location, _ = time.LoadLocation("Europe/Moscow")
var timePoint = time.Date(2020, 1, 28, 16, 20, 39, 50, location)
var events = []event.Event{
	{
		Header:       "Test Event 1",
		StartTime:    timePoint,
		StopTime:     timePoint.Add(time.Hour),
		Description:  func(s string) *string { return &s }("Memory storage test event 1"),
		Owner:        "artem",
		Notification: func(t time.Duration) *time.Duration { return &t }(15 * time.Minute),
	},
	{
		Header:       "Test Event 2",
		StartTime:    timePoint.Add(6 * time.Hour),
		StopTime:     timePoint.Add(6 * time.Hour).Add(30 * time.Minute),
		Description:  func(s string) *string { return &s }("Memory storage test event 2"),
		Owner:        "artem",
		Notification: nil,
	},
	{
		Header:       "Test Event 3",
		StartTime:    timePoint.Add(-time.Hour),
		StopTime:     timePoint.Add(-time.Hour).Add(15 * time.Minute),
		Description:  func(s string) *string { return &s }("Memory storage test event 3"),
		Owner:        "artem",
		Notification: func(t time.Duration) *time.Duration { return &t }(5 * time.Minute),
	},
	{
		Header:       "Test Event 4",
		StartTime:    timePoint.Add(7 * 24 * time.Hour), // Next week
		StopTime:     timePoint.Add(7 * 24 * time.Hour).Add(15 * time.Minute),
		Description:  func(s string) *string { return &s }("Memory storage test event 4"),
		Owner:        "artem",
		Notification: func(t time.Duration) *time.Duration { return &t }(5 * time.Minute),
	},
}

// Factory for memory storage. Build and populate with events
func createMemoryStorage(t *testing.T) *MemoryStorage {
	ms := NewMemoryStorage()
	for _, e := range events {
		if _, err := ms.Add(e); err != nil {
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
		if _, ok := ms.isExistTime(startTime); !ok {
			t.Errorf("event with time %v was not populated", startTime)
		}
	}

	// Test error case
	_, err := ms.Add(events[0])
	assert.NotNil(t, err, "Method should return an error")
}

func TestMemoryStorage_Remove(t *testing.T) {
	var removeTests = struct {
		existed []event.Event
		removed []event.Event
	}{
		[]event.Event{events[0], events[2]},
		[]event.Event{events[1], events[3]},
	}
	ms := createMemoryStorage(t)

	// Remove events from storage
	for _, re := range removeTests.removed {
		id, ok := ms.isExistTime(re.StartTime)
		err := ms.Remove(id)
		assert.Nil(t, err, "Method should return no error")

		// Check if event was removed really
		_, ok = ms.isExistTime(re.StartTime)
		assert.Falsef(t, ok, "event with time %v shouldn't exist in storage", re.StartTime)
	}

	// Check remained events
	for _, ee := range removeTests.existed {
		_, ok := ms.isExistTime(ee.StartTime)
		assert.Truef(t, ok, "event with time %v should exist in storage", ee.StartTime)
	}
}

func TestMemoryStorage_Edit(t *testing.T) {
	ms := createMemoryStorage(t)

	// Modify time for event in storage
	me := events[3]
	me.StartTime = events[0].StartTime.Add(time.Hour)
	id, err := ms.Add(me)
	assert.Nil(t, err, "Method should return no error")
	me.StartTime = me.StartTime.Add(time.Hour)
	err = ms.Edit(id, me)
	assert.Nil(t, err, "Method should return no error")
	me.StartTime = events[0].StartTime
	err = ms.Edit(id, me)
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
