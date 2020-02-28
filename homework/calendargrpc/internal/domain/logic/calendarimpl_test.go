package logic

import (
	"testing"
	"time"

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
	cal := NewCalendar(nil) //storage.NewStorage())
	for _, e := range events {
		if _, err := cal.CreateEvent(e.Title, *e.Description, e.StartTime, e.Duration); err != nil {
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
		events[0].Title, *events[0].Description,
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
		events[0].Title, *events[0].Description,
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
