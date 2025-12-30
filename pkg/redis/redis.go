package redis

import (
	"context"
	"time"

	"github.com/adexcell/delayed-notifier/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr         string        `mapstructure:"addr"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	PoolSize     int           `mapstructure:"pool_size"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	TTL          time.Duration `mapstructure:"ttl"`
}

func New(ctx context.Context, c Config, l *logger.Zerolog) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Addr,
		Password:     c.Password,
		DB:           c.DB,
		MinIdleConns: c.MinIdleConns,
		PoolSize:     c.PoolSize,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		l.Error().Err(err).Msg("redis ping failed")
		return nil, err
	}

	return rdb, nil
}
