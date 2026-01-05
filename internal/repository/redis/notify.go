package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/google/uuid"
	"github.com/wb-go/wbf/redis"
)

type NotifyRedis struct {
	redis *redis.Client
	TTL   time.Duration
}

func NewNotifyRedis(redis *redis.Client, TTL time.Duration) *NotifyRedis {
	return &NotifyRedis{
		redis: redis,
		TTL:   TTL,
	}
}

func (r *NotifyRedis) Get(ctx context.Context, id uuid.UUID) (*domain.Notify, error) {
	key := fmt.Sprintf("notify:%s", id.String())
	res, err := r.redis.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("redis get: %w", err)
	}
	var dto NotifyRedisDTO
	if err := json.Unmarshal([]byte(res), &dto); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return toDomain(dto), nil
}

func (r *NotifyRedis) Set(ctx context.Context, n *domain.Notify) error {
	payload, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("failed to marshal into json: %w", err)
	}
	key := fmt.Sprintf("notify:%s", n.ID.String())
	return r.redis.SetWithExpiration(ctx, key, string(payload),  r.TTL)
}
