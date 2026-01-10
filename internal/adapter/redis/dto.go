package redis

import (
	"encoding/json"
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
)

type NotifyRedisDTO struct {
	ID          string        `json:"notify_id"`
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

func toRedisDTO(n *domain.Notify) ([]byte, error) {
	redistDTO := &NotifyRedisDTO{
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

	payload, err := json.Marshal(redistDTO)
	return payload, err
}

func toDomain(payload string) *domain.Notify {
	var dto NotifyRedisDTO
	json.Unmarshal([]byte(payload), &dto)

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
