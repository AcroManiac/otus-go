package logic

import (
	"encoding/json"
	"io"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/entities"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/interfaces"
	"github.com/pkg/errors"
)

type Scheduler struct {
	collector interfaces.EventsCollector
	wr        io.Writer
}

func NewScheduler(collector interfaces.EventsCollector) interfaces.Scheduler {
	return &Scheduler{collector: collector}
}

func (s *Scheduler) Schedule() error {
	events, err := s.collector.GetEvents()
	if err != nil {
		return errors.Wrap(err, "error collecting events")
	}

	for _, ev := range events {

		// Create notice
		notice := &entities.Notice{
			Title: ev.Title,
			Date:  ev.StartTime,
			Owner: ev.Owner,
		}

		message, err := json.Marshal(notice)
		if err != nil {
			return errors.Wrap(err, "error converting JSON")
		}

		n, err := s.wr.Write(message)
		if err != nil || n != len(message) {
			return errors.Wrap(err, "error sending message")
		}
	}
	return nil
}
