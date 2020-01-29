package storage

// Define error type and interface
type Error string

func (e Error) Error() string {
	return string(e)
}

// Define error messages for storage module
const (
	ErrTimeBusy      Error = "this time is used for other event already"
	ErrNotFoundEvent Error = "event not found"
)
