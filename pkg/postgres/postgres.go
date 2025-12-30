package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/adexcell/delayed-notifier/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	DSN             string        `mapstructure:"dsn"`
	MaxConns        int32         `mapstructure:"max_conns"`
	MinIdleConns    int32         `mapstructure:"min_idle_conns"`
	MaxConnLifetime time.Duration `mapstructure:"max_conn_lifetime"`
	MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time"`
}

func New(ctx context.Context, c Config, l *logger.Zerolog) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(c.DSN)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = c.MaxConns
	poolConfig.MinIdleConns = c.MinIdleConns
	poolConfig.MaxConnLifetime = c.MaxConnLifetime
	poolConfig.MaxConnIdleTime = c.MaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		l.Error().Err(err).Msg("pgxpool init failed")
		return nil, err
	}

	l.Info().
		Str("dsn", c.DSN).
		Int32("max conns", c.MaxConns).
		Dur("max conn lifetime", c.MaxConnLifetime).
		Int32("min iddle conns", c.MinIdleConns).
		Dur("max conn iddle time", c.MaxConnIdleTime).
		Msg("pgxpool created successfully")

	if err := pool.Ping(ctx); err != nil {
		l.Error().Err(err).Msg("conn ping failed")
		return nil, err
	}

	return pool, nil
}
