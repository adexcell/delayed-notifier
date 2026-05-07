package dto

import (
	"time"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
)

// GetNotifyInput represents the input data required to retrieve a specific notification.
type GetNotifyInput struct {
	ID string `binding:"required" json:"id"`
}

// GetNotifyOutput represents the detailed information of a notification.
type GetNotifyOutput struct {
	RecipientEmail string    `json:"recipient_email"`
	Subject        string    `json:"subject"`
	Body           string    `json:"body"`
	RetryCount     int       `json:"retry_count"`
	ScheduledAt    time.Time `json:"scheduled_at"`
	CreatedAt      time.Time `json:"created_at"`
	Status         string    `json:"status"`
	LastError      *string   `json:"last_error,omitempty"`
}

// ToDto converts a domain Notify model to GetNotifyOutput DTO.
func ToDto(n domain.Notify) GetNotifyOutput {
	return GetNotifyOutput{
		RecipientEmail: n.RecipientEmail,
		Subject:        n.Subject,
		Body:           n.Body,
		RetryCount:     n.RetryCount,
		ScheduledAt:    n.ScheduledAt,
		CreatedAt:      n.CreatedAt,
		Status:         n.Status.String(),
		LastError:      n.LastError,
	}
}
