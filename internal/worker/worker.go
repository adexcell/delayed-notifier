package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/wb-go/wbf/zlog"
)

type NotifyConsumer struct {
	postgres domain.NotifyPostgres
	rabbit   domain.QueueProvider
	redis    domain.NotifyRedis
}

func NewNotifyConsumer(postgres domain.NotifyPostgres, rabbit domain.QueueProvider, redis domain.NotifyRedis) *NotifyConsumer {
	return &NotifyConsumer{
		postgres: postgres,
		rabbit:   rabbit,
		redis:    redis,
	}
}

func (c *NotifyConsumer) Handle(ctx context.Context, payload []byte) error {
	var dto NotifyWorkerDTO
	if err := json.Unmarshal(payload, &dto); err != nil {
		zlog.Logger.Error().Err(err).Msg("Consumer: failed to unmarshal")
		return nil
	}

	// TODO: переписать проверку идемпотентности - проверять в redis (в redis надо хранить id обработанных notify)
	// проверка идемпотентности
	currentNotify, err := c.redis.Get(ctx, dto.ID)
	if err != nil {
		currentNotify, err = c.postgres.GetNotifyByID(ctx, dto.ID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				zlog.Logger.Error().Err(err).Any("id", dto.ID).Msgf("Consumer: not found notify %s: %v", dto.ID, err)
				return err
			}
			zlog.Logger.Error().Err(err).Any("id", dto.ID).Msgf("Consumer: failed to fetch current status for %s: %v", dto.ID, err)
			return err
		}

		if currentNotify == nil {
			zlog.Logger.Warn().Any("id", dto.ID).Msgf("Consumer: notify %s not found in DB, skipping", dto.ID)
			return nil
		}
	}

	if currentNotify.Status == domain.StatusSent || currentNotify.Status == domain.StatusCanceled {
		zlog.Logger.Info().Any("id", dto.ID).Msgf("Consumer: notify %s already in final status (%v), skipping", dto.ID, currentNotify.Status)
		return nil
	}
	// -------------------------------------

	// TODO: retryCount добавить
	zlog.Logger.Info().Any("id", dto.ID).Str("Target", dto.Target).Msgf("Consumer: processing notify %s to %s", dto.ID, dto.Target)
	if err := c.Send(ctx, dto); err != nil {
		zlog.Logger.Error().Err(err).Any("id", dto.ID).Msgf("Consumer: send failed for %s: %v", dto.ID, err)
		errStr := err.Error()
		_ = c.postgres.UpdateStatus(ctx, dto.ID, domain.StatusPending, 0, &errStr)
		return nil
	}

	// TODO: retryCount добавить
	if err := c.postgres.UpdateStatus(ctx, dto.ID, domain.StatusSent, 0, nil); err != nil {
		zlog.Logger.Error().Err(err).Any("id", dto.ID).Msg("Consumer: failed to update status to Sent ")
		return err
	}

	zlog.Logger.Info().Any("id", dto.ID).Str("Target", dto.Target).Msg("Consumer: notify sent successfully")
	return nil
}

// TODO: дописать вызов внешнего API  (email, telegram)
func (c *NotifyConsumer) Send(ctx context.Context, dto NotifyWorkerDTO) error {
	if dto.Target == "fail@test.com" {
		return fmt.Errorf("provider connection refused")
	}
	return nil
}
