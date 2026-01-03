package postgres

import (
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/google/uuid"
)

type Status int

const (
	StatusPending   Status = iota // ожидает отправки
	StatusInProcess               // передано в очередь на отправку
	StatusSent                    // отправлено
	StatusFailed                  // ошибка после всех попыток
	StatusCanceled                // отменено пользователем
)

type notifyPostgresDTO struct {
	ID          uuid.UUID `db:"notify_id"`
	Payload     string `db:"payload"`
	Target      string `db:"target"`
	Channel     string `db:"channel"`
	Status      Status `db:"status"`
	ScheduledAt time.Time `db:"scheduled_at"`
}

func toPostgresDTO(n *domain.Notify) *notifyPostgresDTO {
	return &notifyPostgresDTO{
		ID: n.ID,
		Payload: n.Payload,
		Target: n.Target,
		Channel: n.Channel,
		Status: Status(n.Status),
		ScheduledAt: n.ScheduledAt,
	}
}

func toDomain(dto *notifyPostgresDTO) *domain.Notify {
	return &domain.Notify{
		ID: dto.ID,
		Payload: dto.Payload,
		Target: dto.Target,
		Channel: dto.Channel,
		Status: domain.Status(dto.Status),
		ScheduledAt: dto.ScheduledAt,
	}
}
