package usecase

import (
	"context"
	"time"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/adexcell/delayed-notifier/internal/notify/dto"
	"github.com/google/uuid"
)

type Postgres interface {
	CreateNotify(ctx context.Context, notify domain.Notify) error
	GetNotifyStatusByID(ctx context.Context, notifyID uuid.UUID) (domain.Notify, error)
	UpdateNotify(ctx context.Context, notify domain.Notify) error
	DeleteNotify(ctx context.Context, notifyID uuid.UUID) error
}

type Redis interface {
	SetNotifyStatus(ctx context.Context, key, value string, ttl time.Duration) error
	GetNotifyStatus(ctx context.Context, key string) (string, error)
}

type RabbitMQ interface {
	PublishNotify(ctx context.Context, notify domain.Notify) error
}

type NotifyUsecase struct {
	postgres Postgres
	redis    Redis
	rabbit   RabbitMQ
}

func New(postgres Postgres, redis Redis, rabbit RabbitMQ) *NotifyUsecase {
	return &NotifyUsecase{
		postgres: postgres,
		redis:    redis,
		rabbit:   rabbit,
	}
}

func (u *NotifyUsecase) GetNotify(ctx context.Context, input dto.GetNotifyInput) (dto.GetNotifyOutput, error) {
	return dto.GetNotifyOutput{}, nil
}

func (u *NotifyUsecase) DeleteNotify(ctx context.Context, input dto.DeleteNotifyInput) error {
	return nil
}

func (u *NotifyUsecase) UpdateNotify(ctx context.Context, input dto.UpdateNotifyInput) error {
	return nil
}
