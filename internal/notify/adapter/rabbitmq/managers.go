package rabbitmq

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

func (c *Client) manageConsumeChannel(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		conn := c.getConn()
		if conn == nil {
			if !c.sleep(ctx, reInitDelay) {
				return
			}

			continue
		}

		subCh, err := conn.Channel()
		if err != nil {
			log.Warn().Msg("Failed to create sub channel")

			if !c.sleep(ctx, time.Second) {
				return
			}

			continue
		}

		if err := subCh.Qos(prefetchCount, 0, false); err != nil {
			c.clearAndCloseSubChannel(subCh)

			if !c.sleep(ctx, reInitDelay) {
				return
			}

			continue
		}

		c.mu.Lock()

		if c.conn != conn {
			c.mu.Unlock()
			c.clearAndCloseSubChannel(subCh)

			continue
		}

		oldSubCh := c.subCh
		c.subCh = subCh
		c.mu.Unlock()

		if oldSubCh != nil {
			c.subMu.Lock()

			_ = oldSubCh.Close()

			c.subMu.Unlock()
		}

		subClose := make(chan *amqp.Error, 1)
		subCh.NotifyClose(subClose)

		select {
		case <-ctx.Done():
			c.clearAndCloseSubChannel(subCh)

			return

		case err := <-subClose:
			c.clearAndCloseSubChannel(subCh)

			if err != nil {
				log.Warn().Err(err).Msg("[RabbitMQ] Consumer channel closed. Re-init.")
			}

			if !c.sleep(ctx, reInitDelay) {
				return
			}
		}
	}
}

func (c *Client) clearAndCloseSubChannel(subCh *amqp.Channel) {
	if subCh == nil {
		return
	}

	c.mu.Lock()

	if c.subCh == subCh {
		c.subCh = nil
	}

	c.mu.Unlock()

	c.subMu.Lock()

	_ = subCh.Close()

	c.subMu.Unlock()
}

func (c *Client) managePubChannel(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		conn := c.getConn()
		if conn == nil {
			if !c.sleep(ctx, reInitDelay) {
				return
			}

			continue
		}

		pubCh, err := conn.Channel()
		if err != nil {
			log.Warn().Msg("Failed to create pub channel")

			if !c.sleep(ctx, time.Second) {
				return
			}

			continue
		}

		if err := pubCh.Confirm(false); err != nil {
			log.Error().Err(err).Msg("Failed to enable confirms")
			c.clearAndClosePubChannel(pubCh)

			if !c.sleep(ctx, time.Second) {
				return
			}

			continue
		}

		c.mu.Lock()

		if c.conn != conn {
			c.mu.Unlock()
			c.clearAndClosePubChannel(pubCh)

			continue
		}

		oldPubCh := c.pubCh
		c.pubCh = pubCh
		c.mu.Unlock()

		if oldPubCh != nil {
			c.pubMu.Lock()

			_ = oldPubCh.Close()

			c.pubMu.Unlock()
		}

		pubClose := make(chan *amqp.Error, 1)
		pubCh.NotifyClose(pubClose)

		select {
		case <-ctx.Done():
			c.clearAndClosePubChannel(pubCh)

			return

		case err := <-pubClose:
			c.clearAndClosePubChannel(pubCh)

			if err != nil {
				log.Warn().Err(err).Msg("[RabbitMQ] Publisher channel closed. Re-init.")
			}

			if !c.sleep(ctx, reInitDelay) {
				return
			}
		}
	}
}

func (c *Client) clearAndClosePubChannel(pubCh *amqp.Channel) {
	if pubCh == nil {
		return
	}

	c.mu.Lock()

	if c.pubCh == pubCh {
		c.pubCh = nil
	}

	c.mu.Unlock()

	c.pubMu.Lock()

	_ = pubCh.Close()

	c.pubMu.Unlock()
}
