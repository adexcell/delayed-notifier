package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/adexcell/delayed-notifier/pkg/redis"
)

var (
	idempotencyPrefix = "notify:idempotency:"
	ttl               = time.Hour
)

type Redis struct {
	client *redis.Client
}

func New(client *redis.Client) *Redis {
	return &Redis{client: client}
}

func (r *Redis) SetNotifyStatus(ctx context.Context, idempotencyKey, value string) error {
	key := idempotencyPrefix + idempotencyKey

	if err := r.client.Set(ctx, key, value, ttl).Err(); err != nil {
		return fmt.Errorf("redis.Set: %w", err)
	}

	return nil
}

func (r *Redis) GetNotifyStatus(ctx context.Context, idempotencyKey string) (domain.Status, error) {
	key := idempotencyPrefix + idempotencyKey

	val, err := r.client.Get(ctx, key).Result()
	if err == nil {
		return domain.NewStatus(val), err
	}

	return 0, err
}
