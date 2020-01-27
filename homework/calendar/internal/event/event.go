package event

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Id           uuid.UUID
	Header       string
	StartTime    time.Time
	Duration     time.Duration
	Description  *string
	Owner        string
	Notification *time.Duration
}
