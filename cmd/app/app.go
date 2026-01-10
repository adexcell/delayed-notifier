package app

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/adexcell/delayed-notifier/config"
	"github.com/adexcell/delayed-notifier/internal/adapter/postgres"
	"github.com/adexcell/delayed-notifier/internal/adapter/rabbit"
	"github.com/adexcell/delayed-notifier/internal/adapter/redis"
	"github.com/adexcell/delayed-notifier/internal/controller"
	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/adexcell/delayed-notifier/internal/usecase"
	"github.com/adexcell/delayed-notifier/pkg/httpserver"
	"github.com/adexcell/delayed-notifier/pkg/log"
	"github.com/adexcell/delayed-notifier/pkg/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	cfg       *config.Config
	log       log.Log
	router    *router.Router
	server    *http.Server
	scheduler domain.Scheduler
	closers   []func() error
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	log := log.New()

	return &App{
		cfg:    cfg,
		log:    log,
		router: router.New(cfg.Router),
	}, nil
}

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := a.initDependencies(); err != nil {
		return err
	}

	srv := httpserver.New(a.router, a.cfg.HTTPServer, a.log)
	a.addCloser(srv.Close)
	srv.Start()

	go a.scheduler.Run(ctx)

	<-ctx.Done()
	a.log.Info().Msg("Shutting down application...")
	a.shutdown()

	return nil
}

func (a *App) initDependencies() error {
	// Postgres init
	postgres, err := postgres.New(a.cfg.Postgres)
	if err != nil {
		return fmt.Errorf("failed to init Postgres: %w", err)
	}
	a.addCloser(postgres.Close)

	// Redis init
	redis := redis.New(a.cfg.Redis)
	a.addCloser(redis.Close)

	// Rabbit init, declare Queue
	rabbit, err := rabbit.NewRabbitQueueAdapter(a.cfg.Rabbit)
	if err != nil {
		return fmt.Errorf("failed to init Rabbit: %w", err)
	}
	a.addCloser(rabbit.Close)

	// Init Scheduler - producer for notifies
	a.scheduler = usecase.NewScheduler(postgres, rabbit, a.cfg.Scheduler, a.log)

	// Inject dependencies
	notifyUsecase := usecase.New(postgres, redis, rabbit, a.log)
	notifyHandler := controller.NewNotifyHandler(notifyUsecase, a.log)

	// Add static to router, register routers and swagger
	a.router.Static("/static", "./static")
	a.router.StaticFile("/", "./static/index.html")

	notifyHandler.Register(a.router)

	a.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return nil
}

func (a *App) addCloser(closer func() error) {
	a.closers = append(a.closers, closer)
}

func (a *App) shutdown() {
	for i := len(a.closers) - 1; i >= 0; i-- {
		if err := a.closers[i](); err != nil {
			a.log.Error().Err(err).Msg("failed to close resource")
		}
	}
}
