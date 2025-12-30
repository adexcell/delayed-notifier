// pkg/router/health.go
package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	storage *pgxpool.Pool
	cache   *redis.Client
	rabbit  *amqp.Connection
}

func NewHealthHandler(storage *pgxpool.Pool, cache *redis.Client, rabbit *amqp.Connection) *HealthHandler {
	return &HealthHandler{
		storage: storage,
		cache:   cache,
		rabbit:  rabbit,
	}
}

func (h *HealthHandler) Check(c *gin.Context) {
	checks := map[string]bool{}

	// Postgres
	if h.storage != nil {
		checks["postgres"] = h.storage.Ping(c.Request.Context()) == nil
	}

	// Redis
	if h.cache != nil {
		checks["redis"] = h.cache.Ping(c.Request.Context()).Err() == nil
	}

	// RabbitMQ
	if h.rabbit != nil {
		checks["rabbitmq"] = !h.rabbit.IsClosed()
	}

	// Проверка успешна если хотя бы один сервис жив
	allGood := true
	for _, ok := range checks {
		if !ok {
			allGood = false
			break
		}
	}

	if allGood && len(checks) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"checks": checks,
		})
		return
	}

	c.JSON(http.StatusServiceUnavailable, gin.H{
		"status": "error",
		"checks": checks,
	})
}

// Базовый healthcheck без зависимостей
func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "router alive"})
}
