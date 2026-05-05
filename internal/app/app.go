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

func Run(ctx context.Context, c config.Config) (err error) {
	var deps Dependencies

	// ---------Adapters---------
	deps.Postgres, err = postgres.New(ctx, c.Postgres)
	if err != nil {
		return fmt.Errorf("pgxdriver.New: %w", err)
	}

	deps.Redis, err = redis.New(ctx, c.Redis)
	if err != nil {
		return fmt.Errorf("redis.New: %w", err)
	}

	deps.RabbitMQ, err = rabbitmq.New(ctx, c.RabbitMQ)
	if err != nil {
		return fmt.Errorf("rabbitmq.New: %w", err)
	}

	// ---------Controllers---------
	deps.RouterHTTP = ginext.New(c.Router)

	// ---------Metrics---------
	deps.Metrics = metrics.NewHTTPServer()

	// ---------Domains---------
	notify, err := NewNotifyDomain(ctx, deps, c)
	if err != nil {
		return fmt.Errorf("NotifyDomain(deps): %w", err)
	}

	// ---------Start http server---------
	httpserver := httpserver.New(deps.RouterHTTP, c.HTTP)
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
