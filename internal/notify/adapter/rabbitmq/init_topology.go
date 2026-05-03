package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

// initTopology - declare exchanges, queues and bind them.
func (c *Client) initTopology(conn *amqp.Connection) error {
	if conn == nil {
		return errNotConnected
	}

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}
	defer ch.Close()

	// 1. Delayed Exchange
	err = ch.ExchangeDeclare(c.cfg.DelayedExchange, "x-delayed-message", true, false, false, false, amqp.Table{
		"x-delayed-type": "direct",
	})
	if err != nil {
		return fmt.Errorf("declare delayed exchange: %w", err)
	}

	// 2. DLX Exchange
	if err := ch.ExchangeDeclare(c.cfg.DLXExchange, "direct", true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare dlx exchange: %w", err)
	}

	// 3. DLQ Queue
	if _, err := ch.QueueDeclare(c.cfg.DLQQueue, true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare dlq queue: %w", err)
	}

	if err := ch.QueueBind(c.cfg.DLQQueue, c.cfg.DLQRoutingKey, c.cfg.DLXExchange, false, nil); err != nil {
		return fmt.Errorf("bind dlq: %w", err)
	}

	// 4. Main Queue (с привязкой к DLX)
	if _, err := ch.QueueDeclare(c.cfg.MainQueue, true, false, false, false, amqp.Table{
		"x-dead-letter-exchange":    c.cfg.DLXExchange,
		"x-dead-letter-routing-key": c.cfg.DLQRoutingKey,
		"x-queue-type":              "quorum",
	}); err != nil {
		return fmt.Errorf("declare main queue: %w", err)
	}

	if err := ch.QueueBind(c.cfg.MainQueue, c.cfg.RoutingKey, c.cfg.DelayedExchange, false, nil); err != nil {
		return fmt.Errorf("bind main queue: %w", err)
	}

	log.Info().Msg("RabbitMQ topology declared successfully")

	return nil
}
