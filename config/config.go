package config

import (
	"time"

	"github.com/adexcell/delayed-notifier/internal/adapter/sender"
	"github.com/adexcell/delayed-notifier/pkg/httpserver"
	"github.com/adexcell/delayed-notifier/pkg/postgres"
	"github.com/adexcell/delayed-notifier/pkg/rabbit"
	"github.com/adexcell/delayed-notifier/pkg/redis"
	"github.com/adexcell/delayed-notifier/pkg/router"
	"github.com/wb-go/wbf/config"
)

type Config struct {
	App        App                   `mapstructure:"app"`
	HTTPServer httpserver.Config     `mapstructure:"httpserver"`
	Router     router.Config         `mapstructure:"router"`
	Postgres   postgres.Config       `mapstructure:"postgres"`
	Redis      redis.Config          `mapstructure:"redis"`
	Rabbit     rabbit.Config         `mapstructure:"rabbit"`
	Notifier   NotifierConfig        `mapstructure:"notifier"`
	Telegram   sender.TelegramConfig `mapstructure:"telegram"`
	Email sender.EmailConfig    `mapstructure:"email"`
}

type App struct {
	AppName    string `mapstructure:"app_name"`
	AppVersion string `mapstructure:"app_version"`
}

type NotifierConfig struct {
	MaxRetries        int           `mapstructure:"max_retries"`
	VisibilityTimeout time.Duration `mapstructure:"visibility_timeout"`
	Interval          time.Duration `mapstructure:"interval"`
	BatchSize         int           `mapstructure:"batch_size"`
}

func Load() (*Config, error) {
	cfg := config.New()

	cfg.EnableEnv("")

	_ = cfg.LoadEnvFiles(".env")

	if err := cfg.LoadConfigFiles("config/config.yaml"); err != nil {
		return nil, err
	}

	var res Config
	if err := cfg.Unmarshal(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
