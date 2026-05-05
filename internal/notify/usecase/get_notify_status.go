package usecase

import (
	"context"
	"fmt"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/adexcell/delayed-notifier/internal/notify/dto"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (u *NotifyUsecase) GetNotify(ctx context.Context, input dto.GetNotifyInput) (dto.GetNotifyOutput, error) {
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return dto.GetNotifyOutput{}, domain.ErrInvalidID
	}

	status, err := u.redis.GetNotifyStatus(ctx, input.ID)
	if err != nil {
		log.Debug().Err(err).Str("notify id", input.ID).Msg("get notify status cache hit")
	}
	if err == nil {
		return dto.GetNotifyOutput{Status: status}, nil
	}

	status, err = u.postgres.GetNotifyStatusByID(ctx, id)
	if err != nil {
		return dto.GetNotifyOutput{}, fmt.Errorf("failed to get notify status: %w", err)
	}

	// set into redis
	task := domain.NotifyStatusTask{
		Key:   input.ID,
		Value: status.String(),
	}
	u.asyncRedisWriter.Send(ctx, task)

	return dto.GetNotifyOutput{Status: status}, nil
}
