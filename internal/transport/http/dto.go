package http

import (
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/google/uuid"
)

type CreateNotifyRequest struct {
	Payload     string    `json:"payload"`
	Target      string    `json:"target"`
	Channel     string    `json:"channel"`
	ScheduledAt time.Time `json:"scheduled_at"`
}

type NotifyResponse struct {
	ID          uuid.UUID `json:"id"`
	Status      domain.Status `json:"status"`
	ScheduledAt time.Time `json:"scheduled_at"`
	CreatedAt   time.Time `json:"created_at"`
	RetryCount  int `json:"retry_count"`
	LastError   *string `json:"last_error,omitempty"`
}

func toResponse(n *domain.Notify) NotifyResponse {
	return NotifyResponse{
		ID: n.ID,
		Status: n.Status,
		ScheduledAt: n.ScheduledAt,
		CreatedAt: n.CreatedAt,
		RetryCount: n.RetryCount,
		LastError: n.LastError,
	}
}
