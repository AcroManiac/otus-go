package logger

import (
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/entities"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/interfaces"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/monitoring"
	"github.com/prometheus/client_golang/prometheus"
)

// LogSender structure
type LogSender struct {
	stat *prometheus.CounterVec
}

// NewLogSender constructs LogSender structure
func NewLogSender() interfaces.Sender {
	return &LogSender{
		stat: monitoring.NewCounterVec("calendar_sender", "Send", "Send notifications statistics"),
	}
}

// Send function writes data from notice to application log.
// Returns error if any
func (s *LogSender) Send(notice entities.Notice) error {
	Info("Sending notice to event owner", "notice", notice)
	s.stat.WithLabelValues("counter").Inc()
	return nil
}
