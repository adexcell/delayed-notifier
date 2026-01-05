package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/google/uuid"
)

type NotifyUsecase struct {
	postgres domain.NotifyPostgres
	redis    domain.NotifyRedis
	rabbit   domain.QueueProvider
}

func NewNotifyUsecase(
	postgres domain.NotifyPostgres,
	redis domain.NotifyRedis,
	rabbit domain.QueueProvider,
) *NotifyUsecase {
	return &NotifyUsecase{
		postgres: postgres,
		redis:    redis,
		rabbit:   rabbit,
	}
}

func (u *NotifyUsecase) Save(ctx context.Context, n *domain.Notify) (uuid.UUID, error) {
	_, err := u.postgres.GetNotifyByID(ctx, n.ID)
	if err == nil {
		return n.ID, domain.ErrNotifyAlreadyExisis
	}

	if err := u.postgres.Create(ctx, n); err != nil {
		return n.ID, fmt.Errorf("failed to create save notify in db: %w", err)
	}

	if err := u.rabbit.Publish(ctx, n); err != nil {
		return uuid.Nil, fmt.Errorf("failed to publish in broker: %w", err)
	}
	return n.ID, nil
}

func (u *NotifyUsecase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Notify, error) {
	n, err := u.redis.Get(ctx, id)
	if err == nil {
		return n, nil
	}

	n, err = u.postgres.GetNotifyByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFoundNotify) {
			return nil, domain.ErrNotFoundNotify
		}
		return nil, fmt.Errorf("failed to get from db: %w", err)
	}

	if err := u.redis.Set(ctx, n); err != nil {
		return nil, fmt.Errorf("failed to set to redis: %w", err)
	}

	return n, nil
}

func (u *NotifyUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.postgres.DeleteByID(ctx, id)
}

