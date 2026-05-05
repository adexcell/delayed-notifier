package usecase

import (
	"context"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/adexcell/delayed-notifier/internal/notify/dto"
	"github.com/google/uuid"
)

type Postgres interface {
	CreateNotify(ctx context.Context, notify domain.Notify) error
	GetNotifyStatusByID(ctx context.Context, notifyID uuid.UUID) (domain.Status, error)
	DeleteNotify(ctx context.Context, notifyID uuid.UUID) error
}

type Redis interface {
	GetNotifyStatus(ctx context.Context, key string) (domain.Status, error)
}

type AsyncRabbitWriter interface {
	Send(ctx context.Context, delivery dto.Delivery)
}

type AsyncRedisWriter interface {
	Send(ctx context.Context, task domain.NotifyStatusTask)
}

type NotifyUsecase struct {
	postgres          Postgres
	redis             Redis
	asyncRedisWriter  AsyncRedisWriter
	asyncRabbitWriter AsyncRabbitWriter
	notifyStatusCH    chan string
}

func New(postgres Postgres, redis Redis, asyncRedisWriter AsyncRedisWriter, asyncRabbitWriter AsyncRabbitWriter) *NotifyUsecase {
	return &NotifyUsecase{
		postgres:          postgres,
		redis:             redis,
		asyncRedisWriter:  asyncRedisWriter,
		asyncRabbitWriter: asyncRabbitWriter,
		notifyStatusCH:    make(chan string, 1000),
	}
}
