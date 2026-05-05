package dto

import (
	"time"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/google/uuid"
)

type Delivery struct {
	ID             uuid.UUID `json:"id"`
	RecipientEmail string    `json:"recipient_email"`
	Subject        string    `json:"subject"`
	Body           string    `json:"body"`
	ScheduledAt    time.Time `json:"scheduled_at"`
	RetryCount     int       `json:"retry_count"`
	MaxRetries     int       `json:"max_retries"`
	LastError      *string   `json:"last_error,omitempty"`
}

func (d Delivery) ToDomain() domain.Notify {
	return domain.Notify{
		ID:             d.ID,
		RecipientEmail: d.RecipientEmail,
		Subject:        d.Subject,
		Body:           d.Body,
		RetryCount:     d.RetryCount,
		MaxRetries:     d.MaxRetries,
		LastError:      d.LastError,
	}
}

func ToDelivery(n domain.Notify) Delivery {
	return Delivery{
		ID:             n.ID,
		RecipientEmail: n.RecipientEmail,
		Subject:        n.Subject,
		Body:           n.Body,
		RetryCount:     n.RetryCount,
		MaxRetries:     n.MaxRetries,
		LastError:      n.LastError,
	}
}
