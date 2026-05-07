package usecase

import (
	"context"
	"fmt"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/adexcell/delayed-notifier/internal/notify/dto"
	"github.com/adexcell/delayed-notifier/pkg/otel/tracer"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// GetStatus retrieves the current status of a notification, checking cache first then database.
func (u *NotifyUsecase) GetStatus(ctx context.Context, input dto.GetStatusInput) (dto.GetStatusOutput, error) {
	ctx, span := tracer.Start(ctx, "notify usecase GetNotifyStatus")
	defer span.End()

	id, err := uuid.Parse(input.ID)
	if err != nil {
		return dto.GetStatusOutput{}, domain.ErrInvalidID
	}

	status, err := u.redis.GetNotifyStatus(ctx, input.ID)
	if err != nil {
		log.Debug().Err(err).Str("notify id", input.ID).Msg("get notify status cache hit")
	}
	if err == nil {
		return dto.GetStatusOutput{Status: status}, nil
	}

	status, err = u.postgres.GetNotifyStatusByID(ctx, id)
	if err != nil {
		return dto.GetStatusOutput{}, fmt.Errorf("failed to get notify status: %w", err)
	}

	// set into redis
	task := domain.NotifyStatusTask{
		Op:    domain.OpSet,
		Key:   input.ID,
		Value: status.String(),
	}
	u.asyncRedisWriter.Send(ctx, task)

	return dto.GetStatusOutput{Status: status}, nil
}
