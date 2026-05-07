package dto

import (
	"time"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/google/uuid"
)

// Delivery represents the data structure for notification delivery via RabbitMQ.
type Delivery struct {
	ID             uuid.UUID `json:"id"`
	RecipientEmail string    `json:"recipient_email"`
	Subject        string    `json:"subject"`
	Body           string    `json:"body"`
	ScheduledAt    time.Time `json:"scheduled_at"`
	RetryCount     int       `json:"retry_count"`
	MaxRetries     int       `json:"max_retries"`
	LastError      *string   `json:"last_error,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// ToDomain converts a Delivery DTO to a domain Notify model.
func (d Delivery) ToDomain() domain.Notify {
	return domain.Notify{
		ID:             d.ID,
		RecipientEmail: d.RecipientEmail,
		Subject:        d.Subject,
		Body:           d.Body,
		ScheduledAt:    d.ScheduledAt,
		RetryCount:     d.RetryCount,
		MaxRetries:     d.MaxRetries,
		LastError:      d.LastError,
		CreatedAt:      d.CreatedAt,
	}
}

// ToDelivery converts a domain Notify model to a Delivery DTO.
func ToDelivery(n domain.Notify) Delivery {
	return Delivery{
		ID:             n.ID,
		RecipientEmail: n.RecipientEmail,
		Subject:        n.Subject,
		Body:           n.Body,
		ScheduledAt:    n.ScheduledAt,
		RetryCount:     n.RetryCount,
		MaxRetries:     n.MaxRetries,
		LastError:      n.LastError,
		CreatedAt:      n.CreatedAt,
	}
}
