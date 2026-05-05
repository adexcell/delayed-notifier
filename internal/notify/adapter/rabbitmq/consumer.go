package rabbitmq

import (
	"context"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Caller must call Consume again when returned deliveries channel is closed.
func (c *Client) Consume(ctx context.Context) (<-chan amqp.Delivery, error) {
	c.mu.Lock()
	subCh := c.subCh
	queue := c.cfg.MainQueue
	c.mu.Unlock()

	if subCh == nil {
		return nil, errors.New("consumer channel not ready")
	}

	return subCh.ConsumeWithContext(ctx, queue, "", false, false, false, false, nil)
}
