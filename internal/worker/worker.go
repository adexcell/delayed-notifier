package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/adexcell/delayed-notifier/config"
	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/adexcell/delayed-notifier/pkg/log"
)

type NotifyConsumer struct {
	postgres   domain.NotifyPostgres
	rabbit     domain.QueueProvider
	redis      domain.NotifyRedis
	senders    map[string]domain.Sender
	maxRetries int
	log        log.Log
}

func NewNotifyConsumer(
	cfg config.NotifierConfig,
	postgres domain.NotifyPostgres,
	rabbit domain.QueueProvider,
	redis domain.NotifyRedis,
	senders map[string]domain.Sender,
	log log.Log,
) *NotifyConsumer {
	return &NotifyConsumer{
		postgres: postgres,
		rabbit:   rabbit,
		redis:    redis,
		senders:  senders,
		maxRetries: cfg.MaxRetries,
		log:      log,
	}
}

func (c *NotifyConsumer) Handle(ctx context.Context, payload []byte) error {
	var dto NotifyWorkerDTO
	if err := json.Unmarshal(payload, &dto); err != nil {
		c.log.Error().Err(err).Msg("Consumer: failed to unmarshal")
		return nil
	}

	currentNotify, err := c.redis.Get(ctx, dto.ID)
	if err != nil {
		currentNotify, err = c.postgres.GetNotifyByID(ctx, dto.ID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				c.log.Error().Err(err).Any("id", dto.ID).Msgf("Consumer: not found notify %s: %v", dto.ID, err)
				return nil
			}
			c.log.Error().Err(err).Any("id", dto.ID).Msgf("Consumer: failed to fetch current status for %s: %v", dto.ID, err)
			return nil
		}

		if currentNotify == nil {
			c.log.Warn().Any("id", dto.ID).Msgf("Consumer: notify %s not found in DB, skipping", dto.ID)
			return nil
		}
	}

	if currentNotify.Status == domain.StatusSent || currentNotify.Status == domain.StatusCanceled {
		c.log.Info().
			Any("id", dto.ID).
			Msgf("Consumer: notify %s already in final status (%v), skipping", dto.ID, currentNotify.Status)
		return nil
	}

	c.log.Info().
		Any("id", dto.ID).
		Str("Target", dto.Target).
		Msgf("Consumer: processing notify %s to %s", dto.ID, dto.Target)

	if err := c.Send(ctx, dto); err != nil {
		c.log.Error().
			Err(err).
			Any("id", dto.ID).
			Msgf("Consumer: send failed")
		errStr := err.Error()
		dto.RetryCount++
		if dto.RetryCount <= c.maxRetries {
			dto.ScheduledAt = time.Now().Add(time.Duration(dto.RetryCount * dto.RetryCount * int(time.Minute)))
			_ = c.postgres.UpdateStatus(ctx, dto.ID, domain.StatusPending, &dto.ScheduledAt, dto.RetryCount, &errStr)
			return nil
		}

		_ = c.postgres.UpdateStatus(ctx, dto.ID, domain.StatusFailed, nil, dto.RetryCount, &errStr)

		return nil
	}

	if err := c.postgres.UpdateStatus(ctx, dto.ID, domain.StatusSent, nil, dto.RetryCount, nil); err != nil {
		c.log.Error().Err(err).Any("id", dto.ID).Msg("Consumer: failed to update status to Sent ")
		return err
	}

	c.log.Info().Any("id", dto.ID).Str("Target", dto.Target).Msg("Consumer: notify sent successfully")
	return nil
}

func (c *NotifyConsumer) Send(ctx context.Context, dto NotifyWorkerDTO) error {
	sender, ok := c.senders[dto.Channel]
	if !ok {
		return fmt.Errorf("unsupported channel: %s", dto.Channel)
	}
	return sender.Send(ctx, toDomain(&dto))
}
