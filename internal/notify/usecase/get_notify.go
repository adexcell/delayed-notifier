package usecase

import (
	"context"
	"fmt"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/adexcell/delayed-notifier/internal/notify/dto"
	"github.com/adexcell/delayed-notifier/pkg/otel/tracer"
	"github.com/google/uuid"
)

// GetNotify retrieves a notification by its ID.
func (u *NotifyUsecase) GetNotify(ctx context.Context, input dto.GetStatusInput) (dto.GetNotifyOutput, error) {
	ctx, span := tracer.Start(ctx, "notify usecase GetNotify")
	defer span.End()

	id, err := uuid.Parse(input.ID)
	if err != nil {
		return dto.GetNotifyOutput{}, domain.ErrInvalidID
	}

	notify, err := u.postgres.GetNotifyByID(ctx, id)
	if err != nil {
		return dto.GetNotifyOutput{}, fmt.Errorf("failed to get notify status: %w", err)
	}

	output := dto.ToDto(notify)

	return output, nil
}
