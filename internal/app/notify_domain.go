package app

import (
	"context"
	"fmt"

	"github.com/adexcell/delayed-notifier/config"
	"github.com/adexcell/delayed-notifier/internal/notify/adapter/postgres"
	"github.com/adexcell/delayed-notifier/internal/notify/adapter/rabbitmq"
	"github.com/adexcell/delayed-notifier/internal/notify/adapter/redis"
	httprouter "github.com/adexcell/delayed-notifier/internal/notify/controller/http_router"
	"github.com/adexcell/delayed-notifier/internal/notify/usecase"
	"github.com/adexcell/delayed-notifier/internal/notify/worker"
)

type NotifyDomain struct {
	postgres postgres.Postgres
	redis    redis.Redis
	rabbit   *rabbitmq.Client
}

func NewNotifyDomain(ctx context.Context, d Dependencies, c config.Config) (*NotifyDomain, error) {
	domain := &NotifyDomain{}

	if d.Postgres == nil {
		return domain, fmt.Errorf("postgres not init")
	}
	if d.Redis == nil {
		return domain, fmt.Errorf("redis not init")
	}
	if d.RabbitMQ == nil {
		return domain, fmt.Errorf("rabbitmq not init")
	}

	redisWriter := worker.NewAsyncRedisWriter(redis.New(d.Redis))
	redisWriter.Start(ctx)

	rabbitPublisher := worker.NewAsyncRabbitPublisher(d.RabbitMQ, c.RetryStrategy)
	rabbitPublisher.Start(ctx)

	notification := usecase.New(
		postgres.New(d.Postgres),
		redis.New(d.Redis),
		redisWriter,
		rabbitPublisher,
	)

	httprouter.NotifyRouter(d.RouterHTTP, notification, d.Metrics)
	return domain, nil
}

func (d *NotifyDomain) Stop() {

}
