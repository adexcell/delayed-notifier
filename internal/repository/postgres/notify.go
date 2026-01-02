package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/google/uuid"
	"github.com/wb-go/wbf/dbpg"
)



type notifyPostgres struct {
	db *dbpg.DB
}

func NewNotifyPostgres(db *dbpg.DB) *notifyPostgres {
	return &notifyPostgres{db: db}
}

func (p *notifyPostgres) Create(ctx context.Context, n *domain.Notify) error {
	dto := toPostgresDTO(n)
	
	query := `
		insert into notify (notify_id, payload, target, channel, status, scheduled_at, created_at)
		values ($1, $2, $3, $4, $5, $6, $7);`

	_, err := p.db.ExecContext(ctx, query,
		dto.ID, dto.Payload, dto.Target, dto.Channel, dto.Status, dto.ScheduledAt, time.Now())
	return err
}

func (p *notifyPostgres) GetByID(ctx context.Context, id uuid.UUID) (*domain.Notify, error) {
	query := `
		SELECT notify_id, payload, target, channel, status, scheduled_at
		FROM notify WHERE notify_id=$1;`
	var dto notifyPostgresDTO

	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&dto.ID, &dto.Payload, &dto.Target, &dto.Channel, &dto.Status, &dto.ScheduledAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFoundNotify
	}
	return toDomain(&dto), err
}

func (p *notifyPostgres) UpdateStatus(ctx context.Context, id uuid.UUID, status Status, lastErr *string) error {
	query := `
		UPDATE notify
		SET status = $1, last_error = $2, updated_at = NOW()
		WHERE notify_id = $3;`

	res, err := p.db.ExecContext(ctx, query, status, lastErr, id)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFoundNotify
	}
	return nil
}

func (p *notifyPostgres) DeleteByID(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE FROM notify
	WHERE notify_id = $1`

	_, err := p.db.ExecContext(ctx, query, id)

	return err
}

// - обновление статусов у группы notify: StatusPending -> StatusInProcess
func (p *notifyPostgres) LockAndFetchReady(ctx context.Context, limit int) ([]*domain.Notify, error) {
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
		RETURNING notify.notify_id, notify.payload, notify.target, notify.channel,
		notify.scheduled_at;`

	rows, err := p.db.QueryContext(ctx, query, StatusPending, time.Now(), limit, StatusInProcess)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*domain.Notify
	for rows.Next() {
		var dto notifyPostgresDTO
		if err := rows.Scan(&dto.ID, &dto.Payload, &dto.Target, &dto.Channel, &dto.ScheduledAt); err != nil {
			return nil, err
		}
		results = append(results, toDomain(&dto))
	}
	return results, nil
}
