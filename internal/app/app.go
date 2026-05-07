package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/adexcell/delayed-notifier/config"
	"github.com/adexcell/delayed-notifier/internal/notify/adapter/rabbitmq"
	httpserver "github.com/adexcell/delayed-notifier/pkg/http/server"
	"github.com/adexcell/delayed-notifier/pkg/metrics"
	"github.com/adexcell/delayed-notifier/pkg/postgres"
	"github.com/adexcell/delayed-notifier/pkg/redis"
	"github.com/rs/zerolog/log"

	"github.com/wb-go/wbf/ginext"
	kafkav2 "github.com/wb-go/wbf/kafka/kafka-v2"
)

// Dependencies holds all the external components and controllers required by the application.
type Dependencies struct {
	// Adapters
	Postgres      *postgres.Pool
	KafkaProducer *kafkav2.Producer
	Redis         *redis.Client
	RabbitMQ      *rabbitmq.Client

	// Controllers
	RouterHTTP    *ginext.Engine
	KafkaConsumer *kafkav2.Consumer

	Metrics *metrics.HTTPServer
}

// Run initializes and starts the application, handling its lifecycle.
func Run(ctx context.Context, cfg config.Config) (err error) {
	var deps Dependencies

	// ---------Adapters---------
	deps.Postgres, err = postgres.New(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("pgxdriver.New: %w", err)
	}

	deps.Redis, err = redis.New(ctx, cfg.Redis)
	if err != nil {
		return fmt.Errorf("redis.New: %w", err)
	}

	deps.RabbitMQ, err = rabbitmq.New(ctx, cfg.RabbitMQ)
	if err != nil {
		return fmt.Errorf("rabbitmq.New: %w", err)
	}

	// ---------Controllers---------
	deps.RouterHTTP = ginext.New(cfg.Router)

	// ---------Metrics---------
	deps.Metrics = metrics.NewHTTPServer()

	// ---------Domains---------
	notify, err := NewNotifyDomain(deps, cfg)
	if err != nil {
		return fmt.Errorf("NotifyDomain(deps): %w", err)
	}
	notify.Start(ctx)

	// ---------Start http server---------
	httpserver := httpserver.New(deps.RouterHTTP, cfg.HTTP)
	log.Info().Msg("App started!")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig // wait signal

	log.Info().Msg("App got signal to stop")

	notify.Stop()

	// ---------Controllers close---------
	httpserver.Close()

	// ---------Adapters close---------
	deps.RabbitMQ.Close()
	deps.Redis.Close()
	deps.Postgres.Close()

	log.Info().Msg("App stopped!")

	return nil
}
