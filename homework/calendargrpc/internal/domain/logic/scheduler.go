package logic

import (
	"context"
	"encoding/json"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/logger"
	"github.com/gofrs/uuid"
	"io"
	"time"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/entities"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/interfaces"
	"github.com/pkg/errors"
)

type Scheduler struct {
	ctx               context.Context
	collector         interfaces.EventsCollector
	cleaner           interfaces.Cleaner
	wr                io.Writer
	collectorInterval time.Duration
	cleanerInterval   time.Duration
}

func NewScheduler(
	ctx context.Context,
	collector interfaces.EventsCollector,
	cleaner interfaces.Cleaner,
	wr io.Writer,
	collectorInterval time.Duration,
	cleanerInterval time.Duration,
) interfaces.Scheduler {
	return &Scheduler{
		ctx:               ctx,
		collector:         collector,
		cleaner:           cleaner,
		wr:                wr,
		collectorInterval: collectorInterval,
		cleanerInterval:   cleanerInterval,
	}
}

func (s *Scheduler) Start() {
	scheduleTicker := time.NewTicker(s.collectorInterval)
	cleanTicker := time.NewTicker(s.cleanerInterval)

	for {
		select {
		case <-s.ctx.Done():
			logger.Debug("Exit from schedule logic")
			return
		case <-scheduleTicker.C:
			if err := s.Schedule(); err != nil {
				logger.Error("scheduler error", "error", err)
			}
		case <-cleanTicker.C:
			if err := s.Clean(); err != nil {
				logger.Error("cleaner error", "error", err)
			}
		}
	}
}

func (s *Scheduler) Schedule() error {
	if s.collector == nil {
		return errors.New("no events collector")
	}

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
	if s.cleaner == nil {
		return errors.New("no events cleaner")
	}
	if err := s.cleaner.Clean(); err != nil {
		return errors.Wrap(err, "error cleaning retained events")
	}
	return nil
}
