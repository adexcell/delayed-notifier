package domain

import (
	"context"
	"time"
)

type Status int

const (
	StatusPending   Status = iota // 0 - ожидает отправки
	StatusInProcess               // 1 - передано в очередь на отправку
	StatusSent                    // 2 - отправлено
	StatusFailed                  // 3 - ошибка после всех попыток
	StatusCanceled                // 4 - отменено пользователем
)


// Формат ID - uuid.UUID из пакета "github.com/google/uuid" приведенный в формат string
type Notify struct {
	ID          string
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
	GetNotifyByID(ctx context.Context, id string) (*Notify, error)
	UpdateStatus(ctx context.Context, id string, status Status, retryCount int, lastErr *string) error
	DeleteByID(ctx context.Context, id string) error
	LockAndFetchReady(ctx context.Context, limit int) ([]*Notify, error)
	Close() error
}

type NotifyUsecase interface {
	Save(ctx context.Context, n *Notify) (string, error)
	GetByID(ctx context.Context, id string) (*Notify, error)
	Delete(ctx context.Context, id string) error
}

type NotifyRedis interface {
	SetWithExpiration(ctx context.Context, n *Notify) error
	Get(ctx context.Context, id string) (*Notify, error)
	Close() error
}

type Scheduler interface {
	Run(ctx context.Context)
}

type QueueProvider interface {
	Publish(ctx context.Context, n *Notify) error
	Close() error
}
