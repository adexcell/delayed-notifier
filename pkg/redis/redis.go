package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Addr     string `envconfig:"REDIS_ADDR"     required:"true"`
	Password string `envconfig:"REDIS_PASSWORD"`
	DB       int    `envconfig:"REDIS_DB"       default:"0"`
}

type Client struct {
	*redis.Client
}

func New(ctx context.Context, c Config) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       c.DB,
	})

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Warn().Err(err).Msg("Redis connection failed")
	}

	log.Info().Str("redis status", pong).Msg("Connected to Redis")

	return &Client{Client: client}, client.Ping(ctx).Err()
}

func (c *Client) Close() {
	err := c.Client.Close()
	if err != nil {
		log.Error().Err(err).Msg("redis: close")
	}

	log.Info().Msg("redis: closed")
}
