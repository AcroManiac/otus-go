package database

import (
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/entities"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/interfaces"
	"github.com/google/uuid"
	"time"
)

type Storage struct {
	conn Connection
}

func NewDatabaseStorage() interfaces.Storage {
	return &Storage{}
}

// isExistTime is checking event time existence in map.
// Used for unit testing purposes also
func (s Storage) isExistTime(time time.Time) (entities.IdType, bool) {
	//
	return entities.IdType(uuid.UUID{}), false
}

// isExistId is checking event existence in map by Id.
func (s Storage) isExistId(id entities.IdType) bool {
	//
	return false
}

func (s *Storage) Add(ev entities.Event) (entities.IdType, error) {
	// Check if event time is occupied already
	if _, ok := s.isExistTime(ev.StartTime); ok {
		return entities.IdType(uuid.UUID{}), entities.ErrTimeBusy
	}
	id := entities.IdType(uuid.New())

	return id, nil
}

func (s *Storage) Remove(id entities.IdType) error {
	if !s.isExistId(id) {
		return entities.ErrNotFoundEvent
	}

	//

	return nil
}

func (s *Storage) Edit(id entities.IdType, ev entities.Event) error {
	// Check input data for errors
	if !s.isExistId(id) {
		return entities.ErrNotFoundEvent
	}
	if _, ok := s.isExistTime(ev.StartTime); ok {
		return entities.ErrTimeBusy
	}

	//

	return nil
}

func (s Storage) GetEventsByTimePeriod(period entities.TimePeriod, t time.Time) ([]entities.Event, error) {
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

	//
	_ = stopTime

	return selected, nil
}
