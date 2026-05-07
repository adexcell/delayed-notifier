package domain

// Status represents the current state of a notification.
type Status int

const (
	// StatusUnknown represents an undefined notification status.
	StatusUnknown Status = iota
	// StatusPending indicates the notification is waiting to be processed.
	StatusPending
	// StatusProcessing indicates the notification is currently being handled.
	StatusProcessing
	// StatusSent indicates the notification has been successfully delivered.
	StatusSent
	// StatusFailed indicates the notification delivery has failed.
	StatusFailed
	// StatusCancelled indicates the notification has been cancelled.
	StatusCancelled
)

// NewStatus converts a string representation to a Status type.
func NewStatus(s string) Status {
	switch s {
	case "pending":
		return StatusPending
	case "processing":
		return StatusProcessing
	case "sent":
		return StatusSent
	case "failed":
		return StatusFailed
	case "cancelled":
		return StatusCancelled
	default:
		return StatusUnknown
	}
}

// String returns the string representation of the Status.
//
//nolint:exhaustive
func (s Status) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusProcessing:
		return "processing"
	case StatusSent:
		return "sent"
	case StatusFailed:
		return "failed"
	case StatusCancelled:
		return "cancelled"
	default:
		return "unknown"
	}
}
