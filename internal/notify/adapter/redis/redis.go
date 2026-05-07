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

// Redis implements the usecase.Redis interface using Redis.
type Redis struct {
	client *redis.Client
}

// New creates a new instance of the Redis adapter.
func New(client *redis.Client) *Redis {
	return &Redis{client: client}
}

// SetNotifyStatus stores the notification status in Redis with an idempotency key.
func (r *Redis) SetNotifyStatus(ctx context.Context, idempotencyKey, value string) error {
	key := idempotencyPrefix + idempotencyKey

	if err := r.client.Set(ctx, key, value, ttl).Err(); err != nil {
		return fmt.Errorf("redis.Set: %w", err)
	}

	return nil
}

// GetNotifyStatus retrieves the notification status from Redis by its idempotency key.
func (r *Redis) GetNotifyStatus(ctx context.Context, idempotencyKey string) (domain.Status, error) {
	key := idempotencyPrefix + idempotencyKey

	val, err := r.client.Get(ctx, key).Result()
	if err == nil {
		return domain.NewStatus(val), err
	}

	return 0, err
}
