package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateNotifyInput represents the input data required to create a new notification.
type CreateNotifyInput struct {
	RecipientEmail string    `binding:"required,email"      json:"recipient_email"`
	Subject        string    `binding:"required,lte=255"    json:"subject"`
	Body           string    `binding:"required"            json:"body"`
	ScheduledAt    time.Time `binding:"required"            json:"scheduled_at"`
}

// CreateNotifyOutput represents the result of a successful notification creation.
type CreateNotifyOutput struct {
	ID uuid.UUID `json:"id"`
}
