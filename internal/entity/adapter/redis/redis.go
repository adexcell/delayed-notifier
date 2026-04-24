package redis

import "github.com/wb-go/wbf/redis"

type Redis struct {
	client *redis.Client
}

func New(client *redis.Client) *Redis {
	return &Redis{client: client}
}
