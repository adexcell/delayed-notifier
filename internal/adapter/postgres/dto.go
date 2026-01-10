package postgres

import (
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
)

type notifyPostgresDTO struct {
	ID          string        `db:"notify_id"`
	Payload     []byte        `db:"payload"`
	Target      string        `db:"target"`
	Channel     string        `db:"channel"`
	Status      domain.Status `db:"status"`
	ScheduledAt time.Time     `db:"scheduled_at"`
	CreatedAt   time.Time     `db:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at"`
	RetryCount  int           `db:"retry_count"`
	LastError   *string       `db:"last_error"`
}

func toPostgresDTO(n *domain.Notify) *notifyPostgresDTO {
	return &notifyPostgresDTO{
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

func toDomain(dto *notifyPostgresDTO) *domain.Notify {
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
