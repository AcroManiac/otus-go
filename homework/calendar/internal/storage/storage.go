package storage

import (
	"time"

	"github.com/AcroManiac/otus-go/homework/calendar/internal/event"
)

type TimePeriod int

const (
	Day TimePeriod = iota
	Week
	Month
)

type Storage interface {
	Add(event event.Event) error
	Remove(time time.Time) error
	Edit(time time.Time, event event.Event) error
	GetByTimePeriod(time time.Time, period TimePeriod) ([]event.Event, error)
}
