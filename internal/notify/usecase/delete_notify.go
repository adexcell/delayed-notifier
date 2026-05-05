package usecase

import (
	"context"
	"fmt"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/adexcell/delayed-notifier/internal/notify/dto"
	"github.com/google/uuid"
)

func (u *NotifyUsecase) DeleteNotify(ctx context.Context, input dto.DeleteNotifyInput) error {
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return domain.ErrInvalidID
	}

	err = u.postgres.DeleteNotify(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete notify: %w", err)
	}

	return nil
}
