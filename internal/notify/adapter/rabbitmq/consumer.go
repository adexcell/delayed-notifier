package rabbitmq

import (
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Caller must call Consume again when returned deliveries channel is closed.
// for {
// 	deliveries, err := client.Consume()
// 	if err != nil {
// 		time.Sleep(time.Second)
// 		continue
// 	}

//		for msg := range deliveries {
//			// handle
//		}
//	}.
func (c *Client) Consume() (<-chan amqp.Delivery, error) {
	c.mu.Lock()
	subCh := c.subCh
	queue := c.cfg.MainQueue
	c.mu.Unlock()

	if subCh == nil {
		return nil, errors.New("consumer channel not ready")
	}

	return subCh.ConsumeWithContext(c.ctx, queue, "", false, false, false, false, nil)
}
