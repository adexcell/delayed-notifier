package rabbit

import (
	"time"

	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
)

type Config struct {
	URL            string         `mapstructure:"url"`
	ConnectionName string         `mapstructure:"connection_name"`
	ConnectTimeout time.Duration  `mapstructure:"connect_timeout"`
	Heartbeat      time.Duration  `mapstructure:"heartbeat"`
	ReconnectStrat retry.Strategy `mapstructure:"reconnect_strat"`
	ProducingStrat retry.Strategy `mapstructure:"producing_strat"`
	ConsumingStrat retry.Strategy `mapstructure:"consuming_strat"`
}

func NewClient(cfg Config) (*rabbitmq.RabbitClient, error) {
	clientCfg := rabbitmq.ClientConfig{
		URL:            cfg.URL,
		ConnectionName: cfg.ConnectionName,
		ConnectTimeout: cfg.ConnectTimeout,
		Heartbeat:      cfg.Heartbeat,
		ReconnectStrat: cfg.ReconnectStrat,
		ProducingStrat: cfg.ProducingStrat,
		ConsumingStrat: cfg.ConsumingStrat,
	}

	return rabbitmq.NewClient(clientCfg)
}
