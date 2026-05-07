package worker

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/adexcell/delayed-notifier/internal/notify/adapter/rabbitmq"
	"github.com/adexcell/delayed-notifier/internal/notify/dto"
	"github.com/adexcell/delayed-notifier/pkg/retry"
)

// RabbitPublisher defines the interface for publishing notifications to RabbitMQ.
type RabbitPublisher interface {
	PublishNotify(ctx context.Context, delivery dto.Delivery) error
}



// AsyncRabbitPublisher handles asynchronous publishing of notifications to RabbitMQ with retry logic.
type AsyncRabbitPublisher struct {
	rabbit        *rabbitmq.Client
	pubCh         chan dto.Delivery
	wg            sync.WaitGroup
	retryStrategy retry.Config
	draintTimeout time.Duration
}

// NewAsyncRabbitPublisher creates a new instance of AsyncRabbitPublisher.
func NewAsyncRabbitPublisher(rabbit *rabbitmq.Client, retryConfig retry.Config) *AsyncRabbitPublisher {
	return &AsyncRabbitPublisher{
		rabbit:        rabbit,
		pubCh:         make(chan dto.Delivery, 1000),
		retryStrategy: retryConfig,
		draintTimeout: 10 * time.Second,
	}
}

// Start begins the background worker for processing the publishing queue.
func (w *AsyncRabbitPublisher) Start(ctx context.Context) {
	w.wg.Add(1)
	go w.worker(ctx)
	log.Debug().Msg("[AsyncRabbitPublisher] worker start")
}

// Send queues a notification delivery to be published asynchronously.
func (w *AsyncRabbitPublisher) Send(ctx context.Context, delivery dto.Delivery) {
	select {
	case <-ctx.Done():
		log.Error().Err(ctx.Err()).Msg("[AsyncRabbitPublisher] Send cancelled with context.Done()")
		return
	case w.pubCh <- delivery:
	}
}

// Stop gracefully shuts down the publisher by closing the queue and waiting for workers to finish.
func (w *AsyncRabbitPublisher) Stop() {
	close(w.pubCh)
	w.wg.Wait()
}

func (w *AsyncRabbitPublisher) worker(ctx context.Context) {
	defer w.wg.Done()
	for {
		select {
		case <-ctx.Done():
			w.drainRemaining(ctx)
			return
		case delivery, ok := <-w.pubCh:
			if !ok {
				return // channel closed
			}
			w.handleTaskWithRetry(ctx, delivery)
		}
	}
}

func (w *AsyncRabbitPublisher) handleTaskWithRetry(ctx context.Context, delivery dto.Delivery) {

	err := retry.DoWithContext(ctx, w.retryStrategy, func() error {
		return w.rabbit.PublishNotify(ctx, delivery)
	})

	if err != nil {
		log.Warn().Err(err).Msg("[AsyncRabbitPublisher] failed to handleTask")
	}
	log.Debug().Msg("[AsyncRabbitPublisher] success to handleTask")
}

func (w *AsyncRabbitPublisher) drainRemaining(ctx context.Context) {
	timer := time.NewTimer(w.draintTimeout)
	defer timer.Stop()

	for {
		select {
		case task, ok := <-w.pubCh:
			if !ok {
				return
			}
			w.handleTaskWithRetry(ctx, task)
		case <-timer.C:
			return
		}
	}
}
