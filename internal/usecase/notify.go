package usecase

import (
	"context"
	"fmt"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/google/uuid"
)

type NotifyUsecase struct {
	postgres domain.NotifyPostgres
	redis    domain.NotifyRedis
	rabbit   domain.NotifyRabbitMQAdapter
}

func NewNotifyUsecase(
	postgres domain.NotifyPostgres,
	redis domain.NotifyRedis,
	rabbit domain.NotifyRabbitMQAdapter,
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

	return n.ID, nil
}

func (u *NotifyUsecase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Notify, error) {
	// TODO: проверить наличие в redis
	n, err := u.redis.Get(ctx, id)

	// TODO: сходить в бд
}

// func (u *NotifyUsecase) Update(ctx context.Context, n *domain.Notify) error {
	// TODO: сходить в бд
// }

// func (u *NotifyUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	// TODO: сходить в бд
// }

// func (u *NotifyUsecase) GetPending(ctx context.Context, limit int) ([]*domain.Notify, error) {
	// TODO: проверить наличие в redis
	// TODO: сходить в бд
// }
