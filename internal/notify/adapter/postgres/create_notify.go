package postgres

import (
	"context"
	"fmt"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
)

// CreateNotify inserts a new notification into the database.
func (p *Postgres) CreateNotify(ctx context.Context, n domain.Notify) error {
	const sql = `INSERT INTO notifications(id, recipient_email, subject, body, scheduled_at)
				   VALUES($1, $2, $3, $4, $5)`

	_, err := p.pgpool.Exec(ctx, sql, n.ID, n.RecipientEmail, n.Subject, n.Body, n.ScheduledAt)
	if err != nil {
		return fmt.Errorf("CreateNotify: %w", err)
	}
	
	return nil
}
