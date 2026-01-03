package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/adexcell/delayed-notifier/internal/repository/rabbit"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/zlog"
)

type NotifyConsumer struct {
	postgres domain.NotifyPostgres
}

func NewNotifyConsumer(postgres domain.NotifyPostgres) *NotifyConsumer {
	return &NotifyConsumer{postgres: postgres}
}

func (c *NotifyConsumer) Handle(ctx context.Context, msg amqp.Delivery) error {
	var dto rabbit.NotifyRabbitDTO
	if err := json.Unmarshal(msg.Body, &dto); err != nil {
		zlog.Logger.Error().Err(err).Msg("Consumer: failed to unmarshal")
		return nil
	}

	// TODO: переписать проверку идемпотентности - проверять в redis (в redis надо хранить id обработанных notify)
	// проверка идемпотентности
	currentNotify, err := c.postgres.GetByID(ctx, dto.ID)
	if err != nil {
		zlog.Logger.Error().Err(err).Any("id", dto.ID).Msgf("Consumer: failed to fetch current status for %s: %v", dto.ID, err)
		return err
	}

	if currentNotify == nil {
		zlog.Logger.Warn().Any("id", dto.ID).Msgf("Consumer: notify %s not found in DB, skipping", dto.ID)
		return nil
	}

	if currentNotify.Status == domain.StatusSent || currentNotify.Status == domain.StatusCanceled {
		zlog.Logger.Info().Any("id", dto.ID).Msgf("Consumer: notify %s already in final status (%v), skipping", dto.ID, currentNotify.Status)
		return nil
	}
	// -------------------------------------

	zlog.Logger.Info().Any("id", dto.ID).Str("Target", dto.Target).Msgf("Consumer: processing notify %s to %s", dto.ID, dto.Target)
	if err := c.send(ctx, dto); err != nil {
		zlog.Logger.Error().Err(err).Any("id", dto.ID).Msgf("Consumer: send failed for %s: %v", dto.ID, err)
		errStr := err.Error()
		_ = c.postgres.UpdateStatus(ctx, dto.ID, domain.StatusPending, &errStr)
		return nil
	}

	if err := c.postgres.UpdateStatus(ctx, dto.ID, domain.StatusSent, nil); err != nil {
		zlog.Logger.Error().Err(err).Any("id", dto.ID).Msg("Consumer: failed to update status to Sent ")
		return err
	}

	zlog.Logger.Info().Any("id", dto.ID).Str("Target", dto.Target).Msg("Consumer: notify sent successfully")
	return nil
}

// TODO: дописать вызов внешнего API  (email, telegram)
func (c *NotifyConsumer) send(ctx context.Context, dto rabbit.NotifyRabbitDTO) error {
	if dto.Target == "fail@test.com" {
		return fmt.Errorf("provider connection refused")
	}
	return nil
}
