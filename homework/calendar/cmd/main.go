package main

import (
	"log"
	"time"

	"github.com/AcroManiac/otus-go/homework/calendar/internal/calendar"
)

func main() {
	// Create calendar
	var cal calendar.Calendar = calendar.NewCalendar()

	// Create and add event
	if _, err := cal.CreateEvent(time.Now(), time.Now().Add(time.Hour)); err != nil {
		log.Printf("Error adding event: %s", err.Error())
	}

	log.Println("Calendar was created. Bye!")
}
