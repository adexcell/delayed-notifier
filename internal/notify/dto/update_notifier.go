package dto

import (
	"time"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
)

type UpdateNotifyInput struct {
	ID          string         `binding:"required"                      json:"id"`
	Channel     domain.Channel `binding:"required,oneof=email telegram" json:"channel"`
	Recipient   string         `binding:"required"                      json:"recipient"`
	Subject     *string        `json:"subject,omitempty"` // optionally for email, may be nil
	Body        string         `binding:"required"                      json:"body"`
	ScheduledAt time.Time      `binding:"required,gt"                   json:"scheduled_at"`
}
