package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/redis/go-redis/v9"
)

type NotifyRedis struct {
	redis *redis.Client
}

func NewNotifyRedis(redis *redis.Client) *NotifyRedis {
	return &NotifyRedis{redis: redis}
}

func (c *NotifyRedis) Set(ctx context.Context, user *domain.User, ttl time.Duration) error {
	value, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("ошибка при подготовке данных к кешированию: %w", err)
	}
	key := fmt.Sprintf("user:%d", user.ID)
	return c.redis.Set(ctx, key, value, ttl).Err()
}

func (c *NotifyRedis) Get(ctx context.Context, userID int64) (*domain.User, error) {
	key := fmt.Sprintf("user:%d", userID)
	value, err := c.redis.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка при извлечении кешированных данных: %w", err)
	}
	user := domain.User{}
	if err := json.Unmarshal(value, &user); err != nil {
		return nil, fmt.Errorf("ошибка при подготовке извлеченных данных из кеша: %w", err)
	}
	return &user, nil

}
