package dto

import (
	"time"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/google/uuid"
)

type CreateNotifyInput struct {
	Channel     domain.Channel `json:"channel" validate:"required,oneof=email telegram"`
	Recipient   string         `json:"recipient" validate:"required"`
	Subject     *string        `json:"subject,omitempty"`
	Body        string         `json:"body" validate:"required"`
	ScheduledAt time.Time      `json:"scheduled_at" validate:"required,gt=now"`
}

type CreateNotifyOutput struct {
	ID          uuid.UUID      `json:"id"`
	Status      string         `json:"status"`
	ScheduledAt time.Time      `json:"scheduled_at"`
	Channel     domain.Channel `json:"channel"`
	Recipient   string         `json:"recipient"`
}

func(i CreateNotifyInput) ToDomain() domain.Notify {
	return domain.Notify{
		ID:      uuid.New(),
		Channel: i.Channel,
		Recipient: i.Recipient,
		Subject: i.Subject,
		Body: i.Body,
		ScheduledAt: i.ScheduledAt,
	}
}
