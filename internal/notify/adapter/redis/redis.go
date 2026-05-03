package redis

import (
	"context"
	"time"

	"github.com/adexcell/delayed-notifier/pkg/redis"
)

type Redis struct {
	client *redis.Client
}

func New(client *redis.Client) *Redis {
	return &Redis{client: client}
}

func (r *Redis) SetNotifyStatus(ctx context.Context, key, value string, ttl time.Duration) error {
	return nil
}

func (r *Redis) GetNotifyStatus(ctx context.Context, key string) (string, error) {
	return "", nil
}
