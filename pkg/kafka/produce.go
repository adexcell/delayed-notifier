package kafka

import (
	"time"

	"github.com/segmentio/kafka-go"
)

type Config struct {
	Brokers             []string      `mapstructure:"brokers"`
	ClientID            string        `mapstructure:"client_id"`
	ConsumerGroup       string        `mapstructure:"consumer_group"`
	ConsumerWorkerCount int           `mapstructure:"consumer_worker_count"`
	RetryMax            int           `mapstructure:"retry_max"`
	RequiredAcks        int           `mapstructure:"required_acks"`
	MaxWaitTime         time.Duration `mapstructure:"max_wait_time"`
	BatchSize           int           `mapstructure:"batch_size"`
}

func New(c Config) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(c.Brokers...),
		Balancer: &kafka.LeastBytes{},
		BatchSize: c.BatchSize,
	}
}
