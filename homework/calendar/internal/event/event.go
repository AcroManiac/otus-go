package event

import (
	"github.com/google/uuid"
	"time"
)

// Define event id type for using in packages
type IdType uuid.UUID

// Time period for events listing
type TimePeriod int

const (
	Day TimePeriod = iota
	Week
	Month
)

type Event struct {
	Header       string
	StartTime    time.Time
	StopTime     time.Time
	Description  *string
	Owner        string
	Notification *time.Duration
}
