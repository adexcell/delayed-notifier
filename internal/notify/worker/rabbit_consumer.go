package worker

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/adexcell/delayed-notifier/internal/notify/dto"
	"github.com/adexcell/delayed-notifier/internal/notify/metrics"
	"github.com/adexcell/delayed-notifier/pkg/otel/tracer"
	"github.com/adexcell/delayed-notifier/pkg/retry"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

// RabbitConsumer defines the interface for consuming messages from RabbitMQ.
type RabbitConsumer interface {
	Consume(ctx context.Context) (<-chan amqp.Delivery, error)
}

// AsyncRabbitWriter defines the interface for asynchronously writing messages back to RabbitMQ.
type AsyncRabbitWriter interface {
	Send(ctx context.Context, delivery dto.Delivery)
}

// EmailSender defines the interface for sending email notifications.
type EmailSender interface {
	Send(ctx context.Context, n domain.Notify) error
}

// Postgres defines the subset of PostgreSQL operations required by the worker.
type Postgres interface {
	GetNotifyStatusByID(ctx context.Context, notifyID uuid.UUID) (domain.Status, error)
	UpdateNotify(ctx context.Context, notify domain.Notify) error
}

// AsyncRabbitConsumer is a worker that consumes notifications from RabbitMQ and sends them.
type AsyncRabbitConsumer struct {
	rabbit        RabbitConsumer
	rabbitWriter  AsyncRabbitWriter
	postgres      Postgres
	subCh         chan []byte
	wg            sync.WaitGroup
	sender        EmailSender
	retryStrategy retry.Config
	draintTimeout time.Duration
	workerCount   int
	jobsCh        chan amqp.Delivery
}

// NewAsyncRabbitConsumer creates a new instance of AsyncRabbitConsumer.
func NewAsyncRabbitConsumer(
	rabbit RabbitConsumer,
	rabbitWriter AsyncRabbitWriter,
	postgres Postgres,
	retryConfig retry.Config,
	sender EmailSender,
) *AsyncRabbitConsumer {
	return &AsyncRabbitConsumer{
		rabbit:        rabbit,
		rabbitWriter:  rabbitWriter,
		postgres:      postgres,
		sender:        sender,
		retryStrategy: retryConfig,
		draintTimeout: 10 * time.Second,
		workerCount:   10, // Оптимальное количество воркеров
		jobsCh:        make(chan amqp.Delivery, 100),
	}
}

// Start begins the worker pool and dispatcher to handle notifications.
func (w *AsyncRabbitConsumer) Start(ctx context.Context) {
	// Запускаем воркеры
	for i := 0; i < w.workerCount; i++ {
		w.wg.Add(1)
		go w.workerPool(ctx)
	}

	// Запускаем диспетчер, который читает из RabbitMQ и раздает задачи
	w.wg.Add(1)
	go w.dispatcher(ctx)
	log.Debug().Msg("[AsyncRabbitConsumer] worker start")
}

// Stop gracefully shuts down the worker by waiting for all tasks to complete.
func (w *AsyncRabbitConsumer) Stop() {
	w.wg.Wait()
}

func (w *AsyncRabbitConsumer) workerPool(ctx context.Context) {
	defer w.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case d := <-w.jobsCh:
			w.handleDelivery(ctx, d)
		}
	}
}

func (w *AsyncRabbitConsumer) dispatcher(ctx context.Context) {
	defer w.wg.Done()
	delay := time.Second * 1
	for {
		deliveries, err := w.rabbit.Consume(ctx)
		if err != nil {
			log.Warn().Err(err).Msg("[async-rabbit-consumer] failed to get consume channel. Retry")

			select {
			case <-ctx.Done():
				return
			case <-time.After(delay):
				delay = time.Duration(float64(delay) * 2)
				continue
			}
		}

		delay = time.Second * 1 // Сбрасываем задержку при успешном подключении

	ConsumeLoop:
		for {
			select {
			case <-ctx.Done():
				return
			case d, ok := <-deliveries:
				if !ok {
					// Канал закрылся (например, обрыв связи с RabbitMQ)
					// Выходим из внутреннего цикла, чтобы переподключиться
					break ConsumeLoop
				}
				// Отправляем задачу свободному воркеру
				select {
				case <-ctx.Done():
					return
				case w.jobsCh <- d:
				}
			}
		}
	}
}

func (w *AsyncRabbitConsumer) handleDelivery(ctx context.Context, d amqp.Delivery) {
	ctx, span := tracer.Start(ctx, "worker handleDelivery")
	defer span.End()

	var delivery dto.Delivery

	err := json.Unmarshal(d.Body, &delivery)
	if err != nil {
		log.Warn().Err(err).Msg("[rabbit-consumer] json.Unmarshal")
		err = d.Nack(false, false) // -> DLQ
		if err != nil {
			log.Debug().Err(err).Msg("[rabbit] nack failed")
		}
		return
	}

	notify := delivery.ToDomain()
	log.Debug().Msg("[AsyncRabbitConsumer] consumed message from rabbit")

	status, err := w.postgres.GetNotifyStatusByID(ctx, notify.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Warn().Err(err).Msg("[AsyncRabbitConsumer] notify not found in postgres, ack and skip")
			err = d.Ack(false)
			if err != nil {
				log.Debug().Err(err).Msg("[rabbit] ack failed")
			}
			return
		}
		log.Error().Err(err).Msg("[AsyncRabbitConsumer] postgres.GetNotifyStatusByID")
		err = d.Ack(false)
		if err != nil {
			log.Debug().Err(err).Msg("[AsyncRabbitConsumer] ack failed")
		}
		return
	}

	if status == domain.StatusCancelled {
		err = d.Ack(false)
		if err != nil {
			log.Debug().Err(err).Msg("[AsyncRabbitConsumer] ack failed")
		}
		return
	}

	notify.Status = status

	err = w.postgres.UpdateNotify(ctx, notify)
	if err != nil {
		log.Error().Err(err).Msg("[AsyncRabbitConsumer] update status to processing")
		return
	}

	err = w.sender.Send(ctx, notify)
	// Success
	if err == nil {
		err := w.postgres.UpdateNotify(ctx, notify)
		if err != nil {
			log.Error().Err(err).Msg("postgres.UpdateNotify")
		}
		log.Info().Str("id", notify.ID.String()).Msg("notification sent successfully")
		metrics.NotificationsProcessedTotal.Inc()
		err = d.Ack(false)
		if err != nil {
			log.Debug().Err(err).Msg("[rabbit] ack failed")
		}
		return
	}

	log.Error().Err(err).Str("id", notify.ID.String()).Msg("[AsyncRabbitConsumer] failed to send")

	// Failed
	metrics.NotificationsFailedTotal.Inc()
	notify.RetryCount++
	errMsg := err.Error()
	notify.LastError = &errMsg
	

	if notify.RetryCount < notify.MaxRetries {
		log.Warn().
			Err(err).
			Str("id", notify.ID.String()).
			Msg("send failed, scheduling retry")

		nextDelay := CalculateExponentialDelay(notify.RetryCount)
		notify.ScheduledAt = notify.ScheduledAt.Add(nextDelay)
		err := w.postgres.UpdateNotify(ctx, notify)
		if err != nil {
			log.Error().Err(err).Msg("postgres.UpdateNotify")
		}
		delivery = dto.ToDelivery(notify)
		w.rabbitWriter.Send(ctx, delivery)
	}

	err = d.Nack(false, false) // -> DLQ
		if err != nil {
			log.Debug().Err(err).Msg("[rabbit] nack failed")
		}
}

// CalculateExponentialDelay computes the delay for the next retry attempt using exponential backoff.
func CalculateExponentialDelay(retryCount int) time.Duration {
	base := 5 * time.Second
	maxDelay := 1 * time.Hour

	delay := base * (1 << uint(retryCount-1)) // 5s, 10s, 20s, 40s...
	if delay > maxDelay {
		delay = maxDelay
	}
	return delay
}
