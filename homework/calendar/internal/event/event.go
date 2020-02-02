package event

import (
	"time"
)

type Event struct {
	Header       string
	StartTime    time.Time
	StopTime     time.Time
	Description  *string
	Owner        string
	Notification *time.Duration
}
