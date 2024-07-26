package batcher

const (
	// ErrInputChannelClosed should be returned when the batch channel is closed
	ErrInputChannelClosed = Error("batch channel is closed")
	// ErrInputIsNil should be returned when the input data is nil
	ErrInputIsNil = Error("input is nil")
)

// Error is Batcher error
type Error string

// Error returns an error message as a string
func (e Error) Error() string {
	return string(e)
}
