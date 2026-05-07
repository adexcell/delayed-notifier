package postgres

import (
	"github.com/adexcell/delayed-notifier/pkg/postgres"
)

// Postgres implements the usecase.Postgres interface using PostgreSQL.
type Postgres struct {
	pgpool *postgres.Pool
}

// New creates a new instance of the Postgres adapter.
func New(pgpool *postgres.Pool) *Postgres {
	return &Postgres{pgpool: pgpool}
}
