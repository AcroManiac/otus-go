package logic

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"io"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/entities"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/interfaces"
	"github.com/pkg/errors"
)

type Scheduler struct {
	collector interfaces.EventsCollector
	wr        io.Writer
}

func NewScheduler(collector interfaces.EventsCollector, wr io.Writer) interfaces.Scheduler {
	return &Scheduler{collector: collector, wr: wr}
}

func (s *Scheduler) Schedule() error {
	events, err := s.collector.GetEvents()
	if err != nil {
		return errors.Wrap(err, "error collecting events")
	}

	for _, ev := range events {

		// Create notice
		notice := &entities.Notice{
			Id:    uuid.UUID(ev.Id).String(),
			Title: ev.Title,
			Date:  ev.StartTime,
			Owner: ev.Owner,
		}

		// Convert notice to JSON
		message, err := json.Marshal(notice)
		if err != nil {
			return errors.Wrap(err, "error converting JSON")
		}

		// Send message to RabbitMQ broker
		n, err := s.wr.Write(message)
		if err != nil || n != len(message) {
			return errors.Wrap(err, "error sending message")
		}
	}
	return nil
}

// Clean retained events
func (s *Scheduler) Clean() error {
	return nil
}
