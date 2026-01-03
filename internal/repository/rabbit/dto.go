package rabbit

import (
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/google/uuid"
)

type NotifyRabbitDTO struct {
	ID          uuid.UUID `json:"id"`
	Payload     string    `json:"payload"`
	Target      string    `json:"target"`
	Channel     string    `json:"channel"`
	RetryCount  int       `json:"retry_count"`
	ScheduledAt time.Time `json:"scheduled_at"`
}

func toRabbitDTO(n *domain.Notify) *NotifyRabbitDTO {
	return &NotifyRabbitDTO{
		ID:          n.ID,
		Payload:     n.Payload,
		Target:      n.Target,
		Channel:     n.Channel,
		RetryCount:  n.RetryCount,
		ScheduledAt: n.ScheduledAt,
	}
}

func toDomain(dto *NotifyRabbitDTO) *domain.Notify {
	return &domain.Notify{
		ID:          dto.ID,
		Payload:     dto.Payload,
		Target:      dto.Target,
		Channel:     dto.Channel,
		RetryCount:  dto.RetryCount,
		ScheduledAt: dto.ScheduledAt,
	}
}
