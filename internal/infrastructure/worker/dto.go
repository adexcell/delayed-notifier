package worker

import (
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/google/uuid"
)

type NotifyWorkerDTO struct {
	ID          uuid.UUID     `json:"id"`
	Payload     []byte        `json:"payload"`
	Target      string        `json:"target"`
	Channel     string        `json:"channel"`
	Status      domain.Status `json:"status"`
	ScheduledAt time.Time     `json:"scheduled_at"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	RetryCount  int           `json:"retry_count"`
	LastError   *string       `json:"last_error"`
}

func toRabbitDTO(n *domain.Notify) *NotifyWorkerDTO {
	return &NotifyWorkerDTO{
		ID:          n.ID,
		Payload:     n.Payload,
		Target:      n.Target,
		Channel:     n.Channel,
		Status:      n.Status,
		ScheduledAt: n.ScheduledAt,
		CreatedAt:   n.CreatedAt,
		UpdatedAt:   n.UpdatedAt,
		RetryCount:  n.RetryCount,
		LastError:   n.LastError,
	}
}

func toDomain(dto *NotifyWorkerDTO) *domain.Notify {
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
