// Package postgres является оберткой над вспомогательным пакетом wbf/dbpg.
package postgres

import (
	"fmt"
	"time"

	"github.com/wb-go/wbf/dbpg"
)

type DB = dbpg.DB

type Config struct {
	MasterDSN       string        `mapstructure:"master_dsn"`
	SlavesDSN       []string      `mapstructure:"slaves_dsn"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_life_time"`
}

func New(cfg Config) (*DB, error) {
	dbOpts := &dbpg.Options{
		MaxOpenConns:    cfg.MaxOpenConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
	}

	db, err := dbpg.New(cfg.MasterDSN, cfg.SlavesDSN, dbOpts)
	if err != nil {
		return nil, fmt.Errorf("DB connection failed: %w", err)
	}
	if err := db.Master.Ping(); err != nil {
		return nil, fmt.Errorf("DB Ping failed - check your DSN and SSL mode: %w", err)
	}

	return db, nil
}
