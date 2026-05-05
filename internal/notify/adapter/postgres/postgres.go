package postgres

import (
	"github.com/adexcell/delayed-notifier/pkg/postgres"
)

type Postgres struct {
	pgpool *postgres.Pool
}

func New(pgpool *postgres.Pool) *Postgres {
	return &Postgres{pgpool: pgpool}
}
