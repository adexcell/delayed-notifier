package rabbitmq

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// sleep returns true if the timer expires before the client context is canceled.
func (c *Client) sleep(d time.Duration) bool {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-c.ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

// getConn return current conn.
func (c *Client) getConn() *amqp.Connection {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.conn
}
