package logger

import (
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/entities"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/interfaces"
)

// LogSender structure
type LogSender struct{}

// NewLogSender constructs LogSender structure
func NewLogSender() interfaces.Sender {
	return &LogSender{}
}

// Send function writes data from notice to application log.
// Returns error if any
func (s *LogSender) Send(notice entities.Notice) error {
	Info("Sending notice to event owner", "notice", notice)
	return nil
}
