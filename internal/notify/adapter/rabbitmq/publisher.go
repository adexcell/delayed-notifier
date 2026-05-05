package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/adexcell/delayed-notifier/internal/notify/dto"
	amqp "github.com/rabbitmq/amqp091-go"
)

// caller делает retry при получении ошибки.
func (c *Client) PublishNotify(ctx context.Context, notify dto.Delivery) error {
	c.mu.Lock()
	pubCh := c.pubCh
	exchange := c.cfg.DelayedExchange
	routingKey := c.cfg.RoutingKey
	c.mu.Unlock()

	if pubCh == nil {
		return errors.New("publisher channel not ready")
	}

	body, err := json.Marshal(notify)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	delayMs := max(notify.ScheduledAt.UTC().UnixMilli()-time.Now().UTC().UnixMilli(), 0)
	delayMs = min(math.MaxUint32, delayMs)

	pub := amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
		MessageId:    notify.ID.String(),
		Timestamp:    time.Now(),
		Headers: amqp.Table{
			"x-delay":          uint32(delayMs), // max 49.7 days
			"x-retry-count":    uint32(notify.RetryCount),
			"x-correlation-id": notify.ID.String(),
		},
	}

	publishCtx, cancel := context.WithCancel(ctx)
	stop := context.AfterFunc(publishCtx, cancel)

	defer func() {
		stop()
		cancel()
	}()

	c.pubMu.Lock()
	defer c.pubMu.Unlock()

	confirm, err := pubCh.PublishWithDeferredConfirmWithContext(
		publishCtx,
		exchange,
		routingKey,
		false,
		false,
		pub,
	)
	if err != nil {
		return err
	}

	if confirm == nil {
		return errors.New("publisher confirms not enabled")
	}

	ok, err := confirm.WaitContext(publishCtx)
	if err != nil {
		return fmt.Errorf("wait publish confirm: %w", err)
	}

	if !ok {
		return errMsgNacked
	}

	return nil
}
