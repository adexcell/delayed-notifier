package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// Config holds the configuration parameters for the Redis connection.
type Config struct {
	Addr     string `envconfig:"REDIS_ADDR"     required:"true"`
	Password string `envconfig:"REDIS_PASSWORD"`
	DB       int    `envconfig:"REDIS_DB"       default:"0"`
}

// Client wraps the go-redis Client to provide Redis connection functionality.
type Client struct {
	*redis.Client
}

// New creates a new Redis client and checks the connection with Ping.
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

// Close gracefully closes the Redis client connection.
func (c *Client) Close() {
	err := c.Client.Close()
	if err != nil {
		log.Error().Err(err).Msg("redis: close")
	}

	log.Info().Msg("redis: closed")
}
