package interfaces

import "time"

// RetentionPolicy interface for keeping objects in storage
type RetentionPolicy interface {
	// GetDuration returns retention period duration
	GetDuration() time.Duration
}
