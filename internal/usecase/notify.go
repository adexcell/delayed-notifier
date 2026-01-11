package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/adexcell/delayed-notifier/pkg/log"
)

type NotifyUsecase struct {
	log      log.Log
	postgres domain.NotifyPostgres
	redis    domain.NotifyRedis
	rabbit   domain.QueueProvider
}

func New(
	p domain.NotifyPostgres, 
	redis domain.NotifyRedis, 
	rabbit domain.QueueProvider, 
	l log.Log,
	) domain.NotifyUsecase {
	return &NotifyUsecase{
		log:      l,
		postgres: p,
		redis:    redis,
		rabbit:   rabbit,
	}
}

func (u *NotifyUsecase) Save(ctx context.Context, n *domain.Notify) (string, error) {
	_, err := u.postgres.GetNotifyByID(ctx, n.ID)
	if err == nil {
		return n.ID, domain.ErrNotifyAlreadyExisis
	}

	if err := u.postgres.Create(ctx, n); err != nil {
		return n.ID, fmt.Errorf("failed to create save notify in db: %w", err)
	}

	return n.ID, nil
}

func (u *NotifyUsecase) GetByID(ctx context.Context, id string) (*domain.Notify, error) {
	n, err := u.redis.Get(ctx, id)
	if err == nil {
		return n, nil
	}

	n, err = u.postgres.GetNotifyByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get from db: %w", err)
	}

	if err := u.redis.SetWithExpiration(ctx, n); err != nil {
		return nil, fmt.Errorf("failed to set to redis: %w", err)
	}

	return n, nil
}

func (u *NotifyUsecase) Delete(ctx context.Context, id string) error {
	return u.postgres.DeleteByID(ctx, id)
}
