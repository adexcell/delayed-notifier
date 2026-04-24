package app

import (
	"github.com/adexcell/delayed-notifier/internal/entity/adapter/postgres"
	"github.com/adexcell/delayed-notifier/internal/entity/adapter/rabbit"
	"github.com/adexcell/delayed-notifier/internal/entity/adapter/redis"
	httprouter "github.com/adexcell/delayed-notifier/internal/entity/controller/http_router"
	"github.com/adexcell/delayed-notifier/internal/entity/usecase"
)

func EntityDomain(d Dependencies) {
	entityUseCase := usecase.New(
		postgres.New(d.Postgres),
		redis.New(d.Redis),
		rabbit.New(d.RabbitMQ),
	)

	httprouter.EntityRouter(d.RouterHTTP, entityUseCase, d.Metrics)
}
