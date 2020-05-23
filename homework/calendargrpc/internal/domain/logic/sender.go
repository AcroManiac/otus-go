package logic

import (
	"context"
	"encoding/json"
	"io"
	"sync"
	"time"

	"github.com/ahamtat/otus-go/homework/calendargrpc/internal/domain/entities"
	"github.com/ahamtat/otus-go/homework/calendargrpc/internal/domain/interfaces"
	"github.com/ahamtat/otus-go/homework/calendargrpc/internal/infrastructure/logger"
)

type Sender struct {
	rd   io.Reader
	snds []interfaces.Sender
}

func NewSender(rd io.Reader, snds []interfaces.Sender) *Sender {
	return &Sender{
		rd:   rd,
		snds: snds,
	}
}

// Start reading and processing messages from scheduler
func (s *Sender) Start(ctx context.Context) {

	// Create message buffer and mutex to guard it
	var mx sync.Mutex
	buffer := make([]byte, 1024)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			inputNotice := entities.Notice{}

			time.Sleep(100 * time.Millisecond)

			// Read input message
			mx.Lock()
			length, err := s.rd.Read(buffer)
			if err != nil {
				logger.Error("error reading message from broker", "error", err)
				mx.Unlock()
				continue
			}
			if length == 0 && err == nil {
				// Reading channel possibly is to be closed
				mx.Unlock()
				continue
			}

			// Unmarshal input message from JSON
			err = json.Unmarshal(buffer[:length], &inputNotice)
			if err != nil {
				logger.Error("failed unmarshal incoming message", "error", err)
				mx.Unlock()
				continue
			}
			mx.Unlock()

			// Send notice to senders
			for _, sender := range s.snds {
				// Check writer for validness
				if sender == nil {
					logger.Error("notice sender is nil")
					continue
				}
				sender.Send(inputNotice)
			}
		}
	}
}
