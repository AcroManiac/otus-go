package interfaces

// Scheduler interface
type Scheduler interface {
	// Schedule events to be processed by services
	Schedule() error

	// Clean retained events
	Clean() error
}
