package dto

import (
	"time"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/google/uuid"
)

type CreateNotifyInput struct {
	RecipientEmail string    `binding:"required,email" json:"recipient_email"`
	Subject        string    `binding:"lte=255"        json:"subject,omitempty"`
	Body           string    `binding:"required"       json:"body"`
	ScheduledAt    time.Time `binding:"required,gt"    json:"scheduled_at"`
}

type CreateNotifyOutput struct {
	ID uuid.UUID `json:"id"`
}

func (i CreateNotifyInput) ToDomain() domain.Notify {
	return domain.Notify{
		ID:             uuid.New(),
		RecipientEmail: i.RecipientEmail,
		Subject:        i.Subject,
		Body:           i.Body,
		ScheduledAt:    i.ScheduledAt,
	}
}
