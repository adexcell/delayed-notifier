package usecase

import (
	"context"
	"time"

	"github.com/adexcell/delayed-notifier/config"
	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/adexcell/delayed-notifier/pkg/log"
)

type Scheduler struct {
	postgres          domain.NotifyPostgres
	rabbit            domain.QueueProvider
	interval          time.Duration
	batchSize         int
	maxRetries        int
	visibilityTimeout time.Duration
	log               log.Log
}

func NewScheduler(
	postgres domain.NotifyPostgres,
	rabbit domain.QueueProvider,
	cfg config.NotifierConfig,
	log log.Log,
) domain.Scheduler {
	return &Scheduler{
		postgres:          postgres,
		rabbit:            rabbit,
		interval:          cfg.Interval,
		batchSize:         cfg.BatchSize,
		maxRetries:        cfg.MaxRetries,
		visibilityTimeout: cfg.VisibilityTimeout,
		log:               log,
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
	notifies, err := s.postgres.LockAndFetchReady(ctx, s.batchSize, s.visibilityTimeout)
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

			if n.RetryCount < s.maxRetries {
				n.RetryCount += 1
				status = domain.StatusPending
			} else {
				n.RetryCount = 0
				status = domain.StatusFailed
			}

			if err := s.postgres.UpdateStatus(ctx, n.ID, status, &n.ScheduledAt, n.RetryCount, &errStr); err != nil {
				s.log.Error().Err(err).Msg("Scheduler: failed to update status in db")
			}
		}
	}
}
