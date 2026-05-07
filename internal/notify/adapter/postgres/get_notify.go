package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// GetNotifyByID retrieves a notification from the database by its ID.
func (p *Postgres) GetNotifyByID(ctx context.Context, notifyID uuid.UUID) (domain.Notify, error) {
	const sql = `SELECT id, recipient_email, subject, body, scheduled_at, status, created_at 
				 FROM notifications WHERE id = $1`

	var notify domain.Notify
	var status string

	dest := []any{
		&notify.ID,
		&notify.RecipientEmail,
		&notify.Subject,
		&notify.Body,
		&notify.ScheduledAt,
		&status,
		&notify.CreatedAt,
	}

	err := p.pgpool.QueryRow(ctx, sql, notifyID).Scan(dest...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Notify{}, domain.ErrNotFound
		}
		return domain.Notify{}, fmt.Errorf("GetNotifyByID: %w", err)
	}

	notifyStatus := domain.NewStatus(status)
	if notifyStatus == domain.StatusUnknown {
		return domain.Notify{}, fmt.Errorf("got unknown status: %w", err)
	}
	notify.Status = notifyStatus

	return notify, nil
}
