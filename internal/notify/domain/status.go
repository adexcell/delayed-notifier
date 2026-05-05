package domain

type Status int

const (
	StatusUnknown Status = iota
	StatusPending
	StatusProcessing
	StatusSent
	StatusFailed
	StatusCancelled
)

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
