package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/adexcell/delayed-notifier/pkg/redis"
)

const keyPrefix = "notify"

type Redis struct {
	redis      *redis.RDB
	expiration time.Duration
}

func New(cfg redis.Config) domain.NotifyRedis {
	redis := redis.New(cfg)
	expiration := cfg.TTL
	return &Redis{
		redis:      redis,
		expiration: expiration,
	}
}

func (r *Redis) SetWithExpiration(ctx context.Context, n *domain.Notify) error {
	key := fmt.Sprintf("%s:%s", keyPrefix, n.ID)
	value, err := toRedisDTO(n)
	if err != nil {
		return fmt.Errorf("failed to serialize notify into json")
	}
	return r.redis.SetWithExpiration(ctx, key, value, r.expiration)
}

func (r *Redis) Get(ctx context.Context, id string) (*domain.Notify, error) {
	key := fmt.Sprintf("%s:%s", keyPrefix, id)
	payload, err := r.redis.Get(ctx, key)
	if err != nil {
		if err == redis.RedisError {
			return &domain.Notify{}, domain.ErrNotFound
		}
		return &domain.Notify{}, fmt.Errorf("redis error: %w", err)
	}
	n := toDomain(payload)
	return n, nil
}

func (r *Redis) Close() error {
	return r.redis.Close()
}
