package worker

import (
	"context"
	"sync"
	"time"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/rs/zerolog/log"
)

type Redis interface {
	SetNotifyStatus(ctx context.Context, key, value string) error
}

type AsyncRedisWriter struct {
	redis         Redis
	ch            chan domain.NotifyStatusTask
	wg            sync.WaitGroup
	draintTimeout time.Duration
}

func NewAsyncRedisWriter(redis Redis) *AsyncRedisWriter {
	return &AsyncRedisWriter{
		redis:         redis,
		ch:            make(chan domain.NotifyStatusTask, 1000),
		draintTimeout: 10 * time.Second,
	}
}

func (w *AsyncRedisWriter) Start(ctx context.Context) {
	w.wg.Add(1)
	go w.worker(ctx)
}

func (w *AsyncRedisWriter) Send(ctx context.Context, task domain.NotifyStatusTask) {
	select {
	case <-ctx.Done():
		log.Debug().Msgf("AsyncRedisWriter send cancelled with context.Done(): %s", ctx.Err().Error())
		return
	case w.ch <- task:
	}
}

func (w *AsyncRedisWriter) Stop() {
	close(w.ch)
	w.wg.Wait()
}

func (w *AsyncRedisWriter) worker(ctx context.Context) {
	defer w.wg.Done()
	for {
		select {
		case <-ctx.Done():
			w.drainRemaining(ctx)
			return
		case task, ok := <-w.ch:
			if !ok {
				return // channel closed
			}
			w.handleTask(ctx, task)
		}
	}
}

func (w *AsyncRedisWriter) handleTask(ctx context.Context, task domain.NotifyStatusTask) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := w.redis.SetNotifyStatus(ctx, task.Key, task.Value)
	if err != nil {
		log.Warn().Err(err).Str("id", task.Key).Msg("[async-redis-writer] failed to set")
	}
}

func (w *AsyncRedisWriter) drainRemaining(ctx context.Context) {
	timer := time.NewTimer(w.draintTimeout)
	defer timer.Stop()

	for {
		select {
		case task, ok := <-w.ch:
			if !ok {
				return
			}
			w.handleTask(ctx, task)
		case <-timer.C:
			return
		}
	}
}
