package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
)

// UpdateNotify updates the status of a notification in the database.
func (p *Postgres) UpdateNotify(ctx context.Context, n domain.Notify) error {
	const sql = `UPDATE notifications 
				 SET 
				 	status = $1,
					scheduled_at = $2,
					retry_count = $3,
					max_retries = $4,
					last_error = $5,
					updated_at = $6
				 WHERE id = $6`

	_, err := p.pgpool.Exec(
		ctx,
		sql,
		n.Status.String(),
		n.ScheduledAt,
		n.RetryCount,
		n.MaxRetries,
		n.LastError,
		time.Now(),
		n.ID.String(),
	)
	if err != nil {
		return fmt.Errorf("p.pgpool.Exec: %w", err)
	}

	return nil
}
