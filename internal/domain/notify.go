package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Status int

const (
	StatusPending   Status = iota // 0 - ожидает отправки
	StatusInProcess               // 1 - передано в очередь на отправку
	StatusSent                    // 2 - отправлено
	StatusFailed                  // 3 - ошибка после всех попыток
	StatusCanceled                // 4 - отменено пользователем
)

type Notify struct {
	ID          uuid.UUID
	Payload     []byte
	Target      string
	Channel     string
	Status      Status
	ScheduledAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	RetryCount  int
	LastError   *string
}

type NotifyPostgres interface {
	Create(ctx context.Context, n *Notify) error
	GetNotifyByID(ctx context.Context, id uuid.UUID) (*Notify, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status Status, retryCount int, lastErr *string) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
	LockAndFetchReady(ctx context.Context, limit int) ([]*Notify, error)
}

type NotifyUsecase interface {
	Save(ctx context.Context, n *Notify) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Notify, error)
	UpdateNotify(ctx context.Context, n *Notify) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetPending(ctx context.Context, limit int) ([]*Notify, error)
}

type NotifyRedis interface {
	Set(ctx context.Context, n *Notify) error
	Get(ctx context.Context, id uuid.UUID) (*Notify, error)
}

type MessageHandler func(ctx context.Context, payload []byte) error

type QueueProvider interface {
	Publish(ctx context.Context, n *Notify) error
	Consume(ctx context.Context, handler MessageHandler) error
}
