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

type RabbitPublisher interface {
	PublishNotify(ctx context.Context, delivery dto.Delivery) error
}



type AsyncRabbitPublisher struct {
	rabbit        *rabbitmq.Client
	pubCh         chan dto.Delivery
	wg            sync.WaitGroup
	retryStrategy retry.Config
	draintTimeout time.Duration
}

func NewAsyncRabbitPublisher(rabbit *rabbitmq.Client, retryConfig retry.Config) *AsyncRabbitPublisher {
	return &AsyncRabbitPublisher{
		rabbit:        rabbit,
		pubCh:         make(chan dto.Delivery, 1000),
		retryStrategy: retryConfig,
		draintTimeout: 10 * time.Second,
	}
}

func (w *AsyncRabbitPublisher) Start(ctx context.Context) {
	w.wg.Add(1)
	go w.worker(ctx)
}

func (w *AsyncRabbitPublisher) Send(ctx context.Context, delivery dto.Delivery) {
	select {
	case <-ctx.Done():
		log.Error().Err(ctx.Err()).Msg("AsyncRabbitWriter.Send cancelled with context.Done()")
		return
	case w.pubCh <- delivery:
	}
}

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

	retry.DoWithContext(ctx, w.retryStrategy, func() error {
		return w.rabbit.PublishNotify(ctx, delivery)
	})
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
