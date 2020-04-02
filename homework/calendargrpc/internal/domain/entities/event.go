package entities

import (
	"time"

	"github.com/google/uuid"
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

// Structure to represent event entity at interface layer
type Event struct {
	Id          IdType         `json:"id,omitempty"`
	Title       string         `json:"title,omitempty"`
	StartTime   time.Time      `json:"startTime,omitempty"`
	Duration    time.Duration  `json:"duration,omitempty"`
	Description *string        `json:"description,omitempty"`
	Owner       string         `json:"owner,omitempty"`
	Notify      *time.Duration `json:"duration,omitempty"`
}
