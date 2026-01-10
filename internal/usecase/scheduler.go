package usecase

import (
	"context"
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/adexcell/delayed-notifier/pkg/log"
)

const (
	maxRetries int = 3
)

type SchedulerConfig struct {
	Interval  time.Duration `mapstructure:"interval"`
	BatchSize int `mapstructure:"batch_size"`
}

type Scheduler struct {
	postgres  domain.NotifyPostgres
	rabbit    domain.QueueProvider
	interval  time.Duration
	batchSize int // количество одновременно обрабатываемых уведомлений
	log       log.Log
}

func NewScheduler(
	postgres domain.NotifyPostgres,
	rabbit domain.QueueProvider,
	cfg SchedulerConfig,
	log log.Log,
) domain.Scheduler {
	return &Scheduler{
		postgres:  postgres,
		rabbit:    rabbit,
		interval:  cfg.Interval,
		batchSize: cfg.BatchSize,
		log:       log,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.log.Info().Msg("Scheduler started")

	for {
		select {
		case <-ctx.Done():
			s.log.Info().Msg("Scheduler stopped by context")
			return
		case <-ticker.C:
			s.process(ctx)
		}
	}
}

func (s *Scheduler) process(ctx context.Context) {
	// забираем пачку уведомлений из БД (StatusPending -> StatusInProcess)
	notifies, err := s.postgres.LockAndFetchReady(ctx, s.batchSize)
	if err != nil {
		s.log.Error().Err(err).Msg("Scheduler: failed to fetch notifies from db")
	}

	if len(notifies) == 0 {
		return
	}

	for _, n := range notifies {
		if err := s.rabbit.Publish(ctx, n); err != nil {
			s.log.Error().Err(err).Msg("Scheduler: failed to publish notify")
			errStr := err.Error()

			var status domain.Status

			if n.RetryCount <= maxRetries {
				n.RetryCount += 1
				status = domain.StatusPending
			} else {
				n.RetryCount = 0
				status = domain.StatusFailed
			}

			_ = s.postgres.UpdateStatus(ctx, n.ID, status, n.RetryCount, &errStr)
		}
	}
}
