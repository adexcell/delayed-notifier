package usecase

import (
	"context"
	"fmt"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
)

// UpdateNotifyStatus updates the notification.
func (u *NotifyUsecase) UpdateNotifyStatus(ctx context.Context, notify domain.Notify) error {
	// 1. Валидация и обновление в БД
	if err := u.postgres.UpdateNotify(ctx, notify); err != nil {
		return fmt.Errorf("db update failed: %w", err)
	}

	task := domain.NotifyStatusTask{
		Op:  domain.OpDel,
		Key: notify.ID.String(),
	}

	u.asyncRedisWriter.Send(ctx, task)

	return nil
}
