package postgres

import (
	"context"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/adexcell/delayed-notifier/pkg/postgres"
	"github.com/google/uuid"
)

type Postgres struct {
	pgpool *postgres.Pool
}

func New(pgpool *postgres.Pool) *Postgres {
	return &Postgres{pgpool: pgpool}
}

func (p *Postgres) CreateNotify(ctx context.Context, notify domain.Notify) error {
	return nil
}

func (p *Postgres) GetNotifyStatusByID(ctx context.Context, notifyID uuid.UUID) (domain.Notify, error) {
	return domain.Notify{}, nil
}

func (p *Postgres) UpdateNotify(ctx context.Context, notify domain.Notify) error {
	return nil
}

func (p *Postgres) DeleteNotify(ctx context.Context, notifyID uuid.UUID) error {
	return nil
}
