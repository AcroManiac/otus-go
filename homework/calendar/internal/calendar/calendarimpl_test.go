package calendar

import (
	"github.com/ahamtat/otus-go/homework/calendar/internal/event"
	"github.com/ahamtat/otus-go/homework/calendar/internal/storage"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var location, _ = time.LoadLocation("Europe/Moscow")
var timePoint = time.Date(2020, 2, 2, 12, 12, 25, 0, location)
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

// Factory for calendar. Build and populate with events
func createCalendar(t *testing.T) Calendar {
	cal := NewCalendar(storage.NewStorage())
	for _, e := range events {
		if _, err := cal.CreateEvent(e.StartTime, e.StopTime); err != nil {
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
	id, err := cal.CreateEvent(events[0].StartTime.Add(time.Hour), events[0].StopTime.Add(30*time.Minute))
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
	id, err := cal.CreateEvent(events[0].StartTime.Add(time.Hour), events[0].StopTime.Add(30*time.Minute))
	assert.Nil(t, err, "Error should be nil")

	// Delete existing event
	err = cal.DeleteEvent(id)
	assert.Nil(t, err, "Error should be nil")
}

func TestImpl_GetEventsByTimePeriod(t *testing.T) {
	cal := createCalendar(t)
	assert.NotNil(t, cal, "Object should not be nil")

	filtered, err := cal.GetEventsByTimePeriod(events[0].StartTime, event.Week)
	assert.Nil(t, err, "Method should return no error")
	assert.NotEmpty(t, filtered, "Returned slice should not be empty")
}
