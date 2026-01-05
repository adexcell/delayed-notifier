package http

import (
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/google/uuid"
)

type NotifyTransportDTO struct {
	ID          uuid.UUID     `json:"notify_id,omitempty"`
	Payload     []byte        `json:"payload"`
	Target      string        `json:"target"`
	Channel     string        `json:"channel"`
	Status      domain.Status `json:"status"`
	ScheduledAt time.Time     `json:"scheduled_at"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at,omitempty"`
	RetryCount  int           `json:"retry_count,omitempty"`
	LastError   *string       `json:"last_error,omitempty"`
}

type CreateNotifyRequest struct {
	ID          uuid.UUID `json:"id"`
	Payload     string    `json:"payload"`
	Target      string    `json:"target"`
	Channel     string    `json:"channel"`
	ScheduledAt time.Time `json:"scheduled_at"`
}

type NotifyResponse struct {
	ID          uuid.UUID     `json:"id"`
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

func toDomain(dto NotifyTransportDTO) *domain.Notify {
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
