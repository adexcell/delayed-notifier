package controller

import (
	"encoding/json"
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
)

type NotifyControllerDTO struct {
	ID          string          `json:"notify_id,omitempty"`
	Payload     json.RawMessage `json:"payload"`
	Target      string          `json:"target"`
	Channel     string          `json:"channel"`
	Status      domain.Status   `json:"status"`
	ScheduledAt time.Time       `json:"scheduled_at"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at,omitempty"`
	RetryCount  int             `json:"retry_count,omitempty"`
	LastError   *string         `json:"last_error,omitempty"`
}

type CreateNotifyRequest struct {
	ID          string          `json:"id"`
	Payload     json.RawMessage `json:"payload"`
	Target      string          `json:"target"`
	Channel     string          `json:"channel"`
	ScheduledAt time.Time       `json:"scheduled_at"`
}

type NotifyResponse struct {
	ID          string        `json:"id"`
	Status      domain.Status `json:"status"`
	ScheduledAt time.Time     `json:"scheduled_at"`
	CreatedAt   time.Time     `json:"created_at"`
	RetryCount  int           `json:"retry_count"`
	LastError   *string       `json:"last_error,omitempty"`
}

func toResponse(n *domain.Notify) NotifyResponse {
	return NotifyResponse{
		ID:          n.ID,
		Status:      n.Status,
		ScheduledAt: n.ScheduledAt,
		CreatedAt:   n.CreatedAt,
		RetryCount:  n.RetryCount,
		LastError:   n.LastError,
	}
}

func toDomain(dto NotifyControllerDTO) *domain.Notify {
	return &domain.Notify{
		ID:          dto.ID,
		Payload:     dto.Payload,
		Target:      dto.Target,
		Channel:     dto.Channel,
		Status:      dto.Status,
		ScheduledAt: dto.ScheduledAt,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
		RetryCount:  dto.RetryCount,
		LastError:   dto.LastError,
	}
}
