package app

import (
	"github.com/adexcell/delayed-notifier/internal/notify/adapter/postgres"
	"github.com/adexcell/delayed-notifier/internal/notify/adapter/rabbit"
	"github.com/adexcell/delayed-notifier/internal/notify/adapter/redis"
	httprouter "github.com/adexcell/delayed-notifier/internal/notify/controller/http_router"
	"github.com/adexcell/delayed-notifier/internal/notify/usecase"
)

func NotifyDomain(d Dependencies) {
	notification := usecase.New(
		postgres.New(d.Postgres),
		redis.New(d.Redis),
		rabbit.New(d.RabbitMQ),
	)

	httprouter.NotifyRouter(d.RouterHTTP, notification, d.Metrics)
}
