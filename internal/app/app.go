package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/adexcell/delayed-notifier/config"
	httpserver "github.com/adexcell/delayed-notifier/pkg/http/server"
	"github.com/adexcell/delayed-notifier/pkg/metrics"
	"github.com/adexcell/delayed-notifier/pkg/postgres"
	"github.com/rs/zerolog/log"

	"github.com/wb-go/wbf/ginext"
	kafkav2 "github.com/wb-go/wbf/kafka/kafka-v2"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/redis"
)

type Dependencies struct {
	// Adapters
	Postgres      *postgres.Pool
	KafkaProducer *kafkav2.Producer
	Redis         *redis.Client
	RabbitMQ      *rabbitmq.RabbitClient

	// Controllers
	RouterHTTP    *ginext.Engine
	KafkaConsumer *kafkav2.Consumer

	Metrics *metrics.HTTPServer
}

func Run(ctx context.Context, c config.Config) (err error) {
	var deps Dependencies

	// Adapters
	deps.Postgres, err = postgres.New(ctx, c.Postgres)
	if err != nil {
		return fmt.Errorf("pgxdriver.New: %w", err)
	}

	// Controllers
	deps.RouterHTTP = ginext.New(c.Router)

	// Metrics
	deps.Metrics = metrics.NewHTTPServer()

	// Domains
	NotifyDomain(deps)

	// Start http server
	httpserver := httpserver.New(deps.RouterHTTP, c.HTTP)
	log.Info().Msg("App started!")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig // wait signal

	log.Info().Msg("App got signal to stop")

	// Controllers close
	httpserver.Close()

	// Adapters close
	deps.Postgres.Close()

	log.Info().Msg("App stopped!")

	return nil
}
