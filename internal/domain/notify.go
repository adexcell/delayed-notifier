package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Status int

const (
	StatusPending   Status = iota // ожидает отправки
	StatusInProcess               // передано в очередь на отправку
	StatusSent                    // отправлено
	StatusFailed                  // ошибка после всех попыток
	StatusCanceled                // отменено пользователем
)

type Notify struct {
	ID          uuid.UUID
	Payload     string
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
	GetByID(ctx context.Context, id uuid.UUID) (*Notify, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status Status, lastErr *string) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
	LockAndFetchReady(ctx context.Context, limit int) ([]*Notify, error)
}

type QueueProvider interface {
	Publish(ctx context.Context, n *Notify) error
}
