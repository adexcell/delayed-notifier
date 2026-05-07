package config

import (
	"fmt"
	"log"

	"github.com/adexcell/delayed-notifier/internal/notify/adapter/mailer"
	"github.com/adexcell/delayed-notifier/internal/notify/adapter/rabbitmq"
	httpserver "github.com/adexcell/delayed-notifier/pkg/http/server"
	"github.com/adexcell/delayed-notifier/pkg/logger"
	"github.com/adexcell/delayed-notifier/pkg/otel"
	"github.com/adexcell/delayed-notifier/pkg/postgres"
	"github.com/adexcell/delayed-notifier/pkg/redis"
	"github.com/adexcell/delayed-notifier/pkg/retry"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// App holds basic application metadata.
type App struct {
	Name    string `envconfig:"APP_NAME"    required:"true"`
	Version string `envconfig:"APP_VERSION" required:"true"`
}

// Config represents the complete application configuration.
type Config struct {
	App           App
	HTTP          httpserver.Config
	Logger        logger.Config
	OTEL          otel.Config
	Postgres      postgres.Config
	RabbitMQ      rabbitmq.Config
	Redis         redis.Config
	RetryStrategy retry.Config
	Router        string `envconfig:"GIN_MODE"`
	SMTP          mailer.Config
}

// New loads the configuration from environment variables and .env file.
func New() (Config, error) {
	var config Config

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("warning: .env not loaded: %v", err)
	}

	err = envconfig.Process("", &config)
	if err != nil {
		return config, fmt.Errorf("envconfig.Process: %w", err)
	}

	return config, nil
}
