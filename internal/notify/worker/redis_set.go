package worker

import (
	"context"
	"sync"
	"time"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/rs/zerolog/log"
)

// Redis defines the interface for setting notification status in the cache.
type Redis interface {
	SetNotifyStatus(ctx context.Context, key, value string) error
	Del(ctx context.Context, key string) error
}

// AsyncRedisWriter handles asynchronous updates to the notification status cache.
type AsyncRedisWriter struct {
	redis         Redis
	ch            chan domain.NotifyStatusTask
	wg            sync.WaitGroup
	draintTimeout time.Duration
}

// NewAsyncRedisWriter creates a new instance of AsyncRedisWriter.
func NewAsyncRedisWriter(redis Redis) *AsyncRedisWriter {
	return &AsyncRedisWriter{
		redis:         redis,
		ch:            make(chan domain.NotifyStatusTask, 1000),
		draintTimeout: 10 * time.Second,
	}
}

// Start begins the background worker for processing the Redis update queue.
func (w *AsyncRedisWriter) Start(ctx context.Context) {
	w.wg.Add(1)
	go w.worker(ctx)
}

// Send queues a status update task to be processed asynchronously.
func (w *AsyncRedisWriter) Send(ctx context.Context, task domain.NotifyStatusTask) {
	select {
	case <-ctx.Done():
		log.Debug().Msgf("AsyncRedisWriter send cancelled with context.Done(): %s", ctx.Err().Error())
		return
	case w.ch <- task:
	}
}

// Stop gracefully shuts down the writer by closing the queue and waiting for workers to finish.
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

	var err error

	switch task.Op {
	case domain.OpSet:
		err = w.redis.SetNotifyStatus(ctx, task.Key, task.Value)
		if err != nil {
			log.Debug().Err(err).Str("id", task.Key).Msg("[async-redis-writer] failed to set")
		}
	case domain.OpDel:
		err = w.redis.Del(ctx, task.Key)
		if err != nil {
			log.Debug().Err(err).Str("id", task.Key).Msg("[async-redis-writer] failed to delete")
		}
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
