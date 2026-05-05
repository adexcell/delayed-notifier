package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

func (c *Client) connect(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		conn, err := amqp.Dial(c.addr)
		if err != nil {
			log.Warn().Err(err).Msg("[RabbitMQ] Connection failed. Reconnect.")

			if !c.sleep(ctx, reconnectDelay) {
				return
			}

			continue
		}

		log.Info().Msg("[RabbitMQ] Connect successfully.")

		if err := c.initTopology(conn); err != nil {
			c.clearAndCloseConn(conn)

			log.Warn().Err(err).Msg("[RabbitMQ] init topology failed. Reconnect.")

			if !c.sleep(ctx, reconnectDelay) {
				return
			}

			continue
		}

		c.mu.Lock()
		c.conn = conn
		c.mu.Unlock()

		connClose := make(chan *amqp.Error, 1)
		conn.NotifyClose(connClose)

		select {
		case <-ctx.Done():
			c.clearAndCloseConn(conn)

			return

		case err := <-connClose:
			c.clearAndCloseConn(conn)

			if err == nil && ctx.Err() != nil {
				return // graceful shutdown
			}

			log.Warn().Err(err).Msg("[RabbitMQ] Connection lost. Reconnect.")

			if !c.sleep(ctx, reconnectDelay) {
				return
			}
		}
	}
}

// clearAndCloseConn close amqp connection and sets c.conn to nil if conn is actual.
func (c *Client) clearAndCloseConn(conn *amqp.Connection) {
	if conn == nil {
		return
	}

	var (
		pubCh *amqp.Channel
		subCh *amqp.Channel
	)

	c.mu.Lock()

	if c.conn == conn {
		c.conn = nil
		pubCh = c.pubCh
		subCh = c.subCh
		c.pubCh = nil
		c.subCh = nil
	}

	c.mu.Unlock()

	c.clearAndClosePubChannel(pubCh)
	c.clearAndCloseSubChannel(subCh)

	_ = conn.Close()
}
