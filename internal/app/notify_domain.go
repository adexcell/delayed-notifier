package app

import (
	"context"
	"fmt"

	"github.com/adexcell/delayed-notifier/config"
	"github.com/adexcell/delayed-notifier/internal/notify/adapter/mailer"
	"github.com/adexcell/delayed-notifier/internal/notify/adapter/postgres"
	"github.com/adexcell/delayed-notifier/internal/notify/adapter/redis"
	httprouter "github.com/adexcell/delayed-notifier/internal/notify/controller/http_router"
	"github.com/adexcell/delayed-notifier/internal/notify/usecase"
	"github.com/adexcell/delayed-notifier/internal/notify/worker"
)

// NotifyDomain manages the lifecycle and wiring of notification-related components.
type NotifyDomain struct {
	postgres        *postgres.Postgres
	redis           *redis.Redis
	redisWriter     *worker.AsyncRedisWriter
	rabbitPublisher *worker.AsyncRabbitPublisher
	rabbitConsumer  *worker.AsyncRabbitConsumer
}

// NewNotifyDomain wires together the notification domain components and starts the workers.
func NewNotifyDomain(deps Dependencies, c config.Config) (*NotifyDomain, error) {
	domain := &NotifyDomain{}

	if deps.Postgres == nil {
		return domain, fmt.Errorf("postgres not init")
	}
	if deps.Redis == nil {
		return domain, fmt.Errorf("redis not init")
	}
	if deps.RabbitMQ == nil {
		return domain, fmt.Errorf("rabbitmq not init")
	}

	domain.postgres = postgres.New(deps.Postgres)
	domain.redis = redis.New(deps.Redis)
	domain.redisWriter = worker.NewAsyncRedisWriter(redis.New(deps.Redis))
	domain.rabbitPublisher = worker.NewAsyncRabbitPublisher(deps.RabbitMQ, c.RetryStrategy)

	emailSender := mailer.NewEmailSender(c.SMTP)

	domain.rabbitConsumer = worker.NewAsyncRabbitConsumer(
		deps.RabbitMQ,
		domain.rabbitPublisher,
		domain.postgres,
		c.RetryStrategy,
		emailSender,
	)

	notification := usecase.New(
		domain.postgres,
		domain.redis,
		domain.redisWriter,
		domain.rabbitPublisher,
	)

	httprouter.NotifyRouter(deps.RouterHTTP, notification, deps.Metrics)
	return domain, nil
}

// Start the workers within the notification domain.
func (d *NotifyDomain) Start(ctx context.Context) {
	d.redisWriter.Start(ctx)
	d.rabbitPublisher.Start(ctx)
	d.rabbitConsumer.Start(ctx)
}

// Stop gracefully shuts down the components within the notification domain.
func (d *NotifyDomain) Stop() {
	d.redisWriter.Stop()
	d.rabbitPublisher.Stop()
	d.rabbitConsumer.Stop()
}
