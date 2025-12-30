package rabbitmq

import (
	"github.com/adexcell/delayed-notifier/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	URL           string `mapstructure:"url"`
	Exchange      string `mapstructure:"exchange"`
	Kind          string `mapstructure:"kind"`
	DeliveryMode  int    `mapstructure:"delivery_mode"`
	PrefetchCount int    `mapstructure:"prefetch_count"`
}

func New(url string, l logger.Zerolog) (*amqp.Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		l.Error().Err(err).Msg("rabbitmq connection failed")
		return nil, err
	}
	return conn, nil
}
