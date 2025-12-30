package config

import (
	"fmt"
	"strings"

	"github.com/adexcell/delayed-notifier/pkg/httpserver"
	"github.com/adexcell/delayed-notifier/pkg/kafka"
	"github.com/adexcell/delayed-notifier/pkg/logger"
	"github.com/adexcell/delayed-notifier/pkg/postgres"
	"github.com/adexcell/delayed-notifier/pkg/redis"

	_ "github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type App struct {
	Name    string `envconfig:"APP_NAME" required:"true"`
	Version string `envconfig:"APP_VERSION" required:"true"`
}

type Config struct {
	App           App
	HTTP          httpserver.Config
	Logger        logger.Config
	Postgres      postgres.Config
	Redis         redis.Config
	KafkaProducer kafka.Config
}

func InitConfig() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("configs")
	v.AddConfigPath(".")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config.InitConfig: %w", err)
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("config.InitConfig: %w", err)
	}

	return &c, nil
}
