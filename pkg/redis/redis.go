// Package redis является оберткой над вспомогательным пакетом wbf/redis.
package redis

import (
	"time"

	originalRedis "github.com/go-redis/redis/v8"
	"github.com/wb-go/wbf/redis"
)

const RedisError = originalRedis.Nil

type RDB = redis.Client

type Config struct {
	Addr     string        `mapstructure:"addr"`
	Password string        `mapstructure:"password"`
	DB       int           `mapstructure:"db"`
	TTL      time.Duration `mapstructure:"ttl"`
}

func New(cfg Config) *RDB {
	return redis.New(cfg.Addr, cfg.Password, cfg.DB)
}
