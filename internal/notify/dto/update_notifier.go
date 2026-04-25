package dto

import (
	"time"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
)

type UpdateNotifyInput struct {
	ID          string         `json:"id"`
	Channel     domain.Channel `json:"channel" validate:"required,oneof=email telegram"`
	Recipient   string         `json:"recipient" validate:"required"`
	Subject     *string        `json:"subject,omitempty"` // optionally for email, may be nil
	Body        string         `json:"body" validate:"required"`
	ScheduledAt time.Time      `json:"scheduled_at" validate:"required,gt=now"`
	Status      domain.Status  `json:"status"`
	RetryCount  int            `json:"retry_count"`
	MaxRetries  int            `json:"max_retries"`
	LastError   *string        `json:"last_error"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
