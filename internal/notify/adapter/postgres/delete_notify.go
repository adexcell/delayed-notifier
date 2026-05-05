package postgres

import (
	"context"
	"fmt"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/google/uuid"
)

// DeleteNotify deletes notify in db if exists, else returns domain.ErrNotFound
func (p *Postgres) DeleteNotify(ctx context.Context, notifyID uuid.UUID) error {
	const sql = `DELETE FROM notifications WHERE id = $1`

	res, err := p.pgpool.Exec(ctx, sql, notifyID)
	if err != nil {
		return fmt.Errorf("DeleteNotify: %w", err)
	}
	if res.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}
