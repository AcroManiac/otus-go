package main

import (
	"github.com/AcroManiac/otus-go/homework/calendar/internal/event"
	"github.com/AcroManiac/otus-go/homework/calendar/internal/storage"
	"github.com/google/uuid"
	"log"
	"time"
)

type Calendar struct {
	Storage storage.Storage
}

func main() {
	// Create calendar
	var calendar Calendar = Calendar{Storage: storage.New()}

	// Create and add event
	testEvent := event.Event{
		Id:           uuid.New(),
		Header:       "test event",
		StartTime:    time.Now(),
		Duration:     time.Hour,
		Description:  func(s string) *string { return &s }("Event for calendar testing"),
		Owner:        "artem",
		Notification: func(t time.Duration) *time.Duration { return &t }(15 * time.Minute),
	}
	if err := calendar.Storage.Add(testEvent); err != nil {
		log.Printf("Error adding event: %s", err.Error())
	}

	log.Println("Calendar was created. Bye!")
}
