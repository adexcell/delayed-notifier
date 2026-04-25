package usecase

import (
	"context"

	"github.com/adexcell/delayed-notifier/internal/notify/dto"
	"github.com/adexcell/delayed-notifier/pkg/otel/tracer"
)

type Postgres interface{}

type Redis interface{}

type RabbitMQ interface{}

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

func (u *NotifyUsecase) CreateNotify(ctx context.Context, input dto.CreateNotifyInput) (dto.CreateNotifyOutput, error) {
	ctx, span := tracer.Start(ctx, "notify usecase CreateNotify")
	defer span.End()

	return dto.CreateNotifyOutput{}, nil
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
