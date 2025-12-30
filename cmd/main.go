package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/adexcell/delayed-notifier/config"
	ginrouter "github.com/adexcell/delayed-notifier/internal/controller/http/gin_router"
	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/adexcell/delayed-notifier/pkg/auth"
	"github.com/adexcell/delayed-notifier/pkg/httpserver"
	"github.com/adexcell/delayed-notifier/pkg/kafka"
	"github.com/adexcell/delayed-notifier/pkg/logger"
	"github.com/adexcell/delayed-notifier/pkg/postgres"
	"github.com/adexcell/delayed-notifier/pkg/redis"
	"github.com/adexcell/delayed-notifier/pkg/router"
)

func main() {
	c, err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	logger.Init(c.Logger)

	err = AppRun(context.Background(), c, &logger.Logger)
	if err != nil {
		panic(err)
	}
}

func AppRun(ctx context.Context, c *config.Config, l *logger.Zerolog) error {
	pgPool, err := postgres.New(ctx, c.Postgres, &logger.Logger)
	if err != nil {
		return fmt.Errorf("postgres.New: %w", err)
	}

	redisClient, err := redis.New(ctx, c.Redis, &logger.Logger)
	if err != nil {
		return fmt.Errorf("redis.New: %w", err)
	}

	kafkaProducer := kafka.New(c.KafkaProducer)

	us := domain.NewUser()
	ns := domain.Notify
	tm := auth.TokenManager
	r := ginrouter.NewHandler(us, ns, tm, l)
	httpServer := httpserver.New(r, c.HTTP, &logger.Logger)
	if err := httpServer.Start(); err != nil {
		panic(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig // Ctrl+C or SIGTERM

	// createProfileConsumer.Close()
	httpServer.Stop()
	redisClient.Close()
	kafkaProducer.Close()
	pgPool.Close()
    return nil
}
