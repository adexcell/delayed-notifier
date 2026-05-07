package usecase

import (
	"context"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/adexcell/delayed-notifier/internal/notify/dto"
	"github.com/google/uuid"
)

// Postgres defines the interface for PostgreSQL database operations related to notifications.
type Postgres interface {
	CreateNotify(ctx context.Context, notify domain.Notify) error
	GetNotifyStatusByID(ctx context.Context, notifyID uuid.UUID) (domain.Status, error)
	GetNotifyByID(ctx context.Context, notifyID uuid.UUID) (domain.Notify, error)
	DeleteNotify(ctx context.Context, notifyID uuid.UUID) error
}

// Redis defines the interface for Redis cache operations.
type Redis interface {
	GetNotifyStatus(ctx context.Context, key string) (domain.Status, error)
}

// AsyncRabbitWriter defines the interface for asynchronous writing to RabbitMQ.
type AsyncRabbitWriter interface {
	Send(ctx context.Context, delivery dto.Delivery)
}

// AsyncRedisWriter defines the interface for asynchronous writing to Redis.
type AsyncRedisWriter interface {
	Send(ctx context.Context, task domain.NotifyStatusTask)
}

// NotifyUsecase implements the business logic for notification management.
type NotifyUsecase struct {
	postgres          Postgres
	redis             Redis
	asyncRedisWriter  AsyncRedisWriter
	asyncRabbitWriter AsyncRabbitWriter
	notifyStatusCH    chan string
}

// New creates a new instance of NotifyUsecase.
func New(postgres Postgres, redis Redis, asyncRedisWriter AsyncRedisWriter, asyncRabbitWriter AsyncRabbitWriter) *NotifyUsecase {
	return &NotifyUsecase{
		postgres:          postgres,
		redis:             redis,
		asyncRedisWriter:  asyncRedisWriter,
		asyncRabbitWriter: asyncRabbitWriter,
		notifyStatusCH:    make(chan string, 1000),
	}
}
