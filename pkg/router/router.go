package router

import (
	"github.com/adexcell/delayed-notifier/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type GinRouter struct {
	Router *gin.Engine
}

func New(l *logger.Zerolog) *GinRouter {
	r := &GinRouter{Router: gin.New()}

	r.Router.Use(gin.Recovery(), Logger(l))

	r.Router.GET("/health", healthCheckHandler)
	return r
}

func (r *GinRouter) HealthCheck(storage *pgxpool.Pool, cache *redis.Client, rabbit *amqp.Connection) {
	handler := NewHealthHandler(storage, cache, rabbit)
	r.Router.GET("/health", handler.Check)
}
