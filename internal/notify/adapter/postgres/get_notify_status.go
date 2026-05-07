package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// GetNotifyStatusByID retrieves the status of a notification from the database by its ID.
func (p *Postgres) GetNotifyStatusByID(ctx context.Context, notifyID uuid.UUID) (domain.Status, error) {
	const sql = `SELECT status FROM notifications WHERE id = $1`

	var status string

	err := p.pgpool.QueryRow(ctx, sql, notifyID).Scan(&status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, domain.ErrNotFound
		}
		return 0, fmt.Errorf("GetNotifyStatusByID: %w", err)
	}

	notifyStatus := domain.NewStatus(status)
	if notifyStatus == domain.StatusUnknown {
		return 0, fmt.Errorf("got unknown status: %w", err)
	}

	return notifyStatus, nil
}
