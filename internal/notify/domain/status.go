package domain

type Status int

const (
	Unknown Status = iota
	Pending
	Processing
	Sent
	Failed
	Cancelled
)

func NewStatus(s string) Status {
	switch s {
	case "pending":
		return Pending
	case "processing":
		return Processing
	case "sent":
		return Sent
	case "failed":
		return Failed
	case "cancelled":
		return Cancelled
	default:
		return Unknown
	}
}

//nolint:exhaustive
func (s Status) String() string {
	switch s {
	case Pending:
		return "pending"
	case Processing:
		return "processing"
	case Sent:
		return "sent"
	case Failed:
		return "failed"
	case Cancelled:
		return "cancelled"
	default:
		return "unknown"
	}
}
