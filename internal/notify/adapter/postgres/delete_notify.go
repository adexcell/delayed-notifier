package postgres

import (
	"context"
	"fmt"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/google/uuid"
)

// DeleteNotify removes a notification from the database by its ID.
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
