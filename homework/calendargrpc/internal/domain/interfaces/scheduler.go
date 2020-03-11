package interfaces

// Scheduler interface
type Scheduler interface {
	// Function manages events to be processed by services.
	// Returns error if failed
	Schedule() error
}
