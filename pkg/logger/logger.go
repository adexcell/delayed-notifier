package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config holds the configuration parameters for the logger.
type Config struct {
	AppName       string `envconfig:"APP_NAME"    required:"true"`
	AppVersion    string `envconfig:"APP_VERSION" required:"true"`
	Level         string `default:"debug"         envconfig:"LOGGER_LEVEL"`
	PrettyConsole bool   `default:"false"         envconfig:"LOGGER_PRETTY_CONSOLE"`
}

// Init initializes the global logger with the provided configuration.
func Init(c Config) {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	level, err := zerolog.ParseLevel(c.Level)
	if err != nil {
		level = zerolog.InfoLevel
		log.Warn().Err(err).Msgf("Invalid log level '%s', using %s", c.Level, level)
	}
	zerolog.SetGlobalLevel(level)

	if c.PrettyConsole {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})
		log.Logger = log.With().Caller().Logger()
	} else {
		log.Logger = log.With().
			Str("app_name", c.AppName).
			Str("app_version", c.AppVersion).
			Logger()
	}

	log.Info().Msg("Logger initialized")
}
