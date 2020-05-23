package interfaces

import "github.com/ahamtat/otus-go/homework/calendargrpc/internal/domain/entities"

// Sender interface
type Sender interface {
	// Send function writes data from notice to output.
	// Returns error if any
	Send(notice entities.Notice) error
}
