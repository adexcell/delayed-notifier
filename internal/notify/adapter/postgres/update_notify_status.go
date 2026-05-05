package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (p *Postgres) UpdateNotifyStatus(ctx context.Context, notifyID uuid.UUID, status string) error {
	const sql = `UPDATE notifications SET status = $1 WHERE id = $2`

	_, err := p.pgpool.Exec(ctx, sql, status, notifyID)
	if err != nil {
		return fmt.Errorf("p.pgpool.Exec: %w", err)
	}
	
	return nil
}
