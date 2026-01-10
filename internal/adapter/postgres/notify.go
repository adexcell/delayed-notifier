package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/adexcell/delayed-notifier/pkg/postgres"
)

type Postgres struct {
	db *postgres.DB
}

func New(cfg postgres.Config) (domain.NotifyPostgres, error) {
	db, err := postgres.New(cfg)
	return &Postgres{db: db}, err
}

func (p *Postgres) Create(ctx context.Context, n *domain.Notify) error {
	dto := toPostgresDTO(n)

	query := `
		INSERT INTO notify (notify_id, payload, target, channel, status, scheduled_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7);`

	_, err := p.db.ExecContext(ctx, query,
		dto.ID, dto.Payload, dto.Target, dto.Channel, dto.Status, dto.ScheduledAt, time.Now().UTC())
	return err
}

func (p *Postgres) GetNotifyByID(ctx context.Context, id string) (*domain.Notify, error) {
	query := `
		SELECT notify_id, payload, target, channel, status, scheduled_at
		FROM notify WHERE notify_id=$1;`
	var dto notifyPostgresDTO

	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&dto.ID, &dto.Payload, &dto.Target, &dto.Channel, &dto.Status, &dto.ScheduledAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return toDomain(&dto), err
}

func (p *Postgres) UpdateStatus(
	ctx context.Context,
	id string,
	status domain.Status,
	retryCount int,
	lastErr *string,
) error {
	query := `
		UPDATE notify
		SET status      = $2, 
			retry_count = $3,
			last_error  = $4, 
			updated_at  = NOW()
		WHERE notify_id = $1;`

	res, err := p.db.ExecContext(ctx, query, id, status, retryCount, lastErr)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (p *Postgres) DeleteByID(ctx context.Context, id string) error {
	query := `
	DELETE FROM notify
	WHERE notify_id = $1`

	_, err := p.db.ExecContext(ctx, query, id)

	return err
}

// - обновление статусов у группы notify: StatusPending -> StatusInProcess
func (p *Postgres) LockAndFetchReady(ctx context.Context, limit int) ([]*domain.Notify, error) {
	query := `
		WITH selected AS (
			SELECT notify_id FROM notify
			WHERE status = $1 AND scheduled_at <= $2
			ORDER BY scheduled_at ASC
			LIMIT $3
			FOR UPDATE SKIP LOCKED
		)
		UPDATE notify
		SET status = $4, updated_at = NOW()
		FROM selected
		WHERE notify.notify_id = selected.notify_id
		RETURNING   notify.notify_id, 
					notify.payload, 
					notify.target, 
					notify.channel,
					notify.status, 
					notify.scheduled_at, 
					notify.created_at, 
					notify.updated_at,
					notify.retry_count, 
					notify.last_error;`

	rows, err := p.db.QueryContext(ctx, query, domain.StatusPending, time.Now().UTC(), limit, domain.StatusInProcess)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*domain.Notify
	for rows.Next() {
		var dto notifyPostgresDTO
		if err := rows.Scan(
			&dto.ID,
			&dto.Payload,
			&dto.Target,
			&dto.Channel,
			&dto.Status,
			&dto.ScheduledAt,
			&dto.CreatedAt,
			&dto.UpdatedAt,
			&dto.RetryCount,
			&dto.LastError,
		); err != nil {
			return nil, err
		}
		results = append(results, toDomain(&dto))
	}
	return results, nil
}

func (p *Postgres) Close() error {
	return p.db.Master.Close()
}
