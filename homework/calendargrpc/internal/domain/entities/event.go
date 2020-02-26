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
	Id          IdType
	Title       string
	StartTime   time.Time
	Duration    time.Duration
	Description *string
	Owner       string
	Notify      *time.Duration
}
