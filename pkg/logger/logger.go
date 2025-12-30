package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

type Config struct {
	Level      string `mapstructure:"level" default:"info"`
	JSONFormat bool   `mapstructure:"json_format"`
}

type Zerolog = zerolog.Logger

var Logger Zerolog

func Init(c Config) {
	l, err := zerolog.ParseLevel(c.Level)
	if err != nil {
		l = zerolog.InfoLevel
	}

	var output io.Writer

	if c.JSONFormat {
		output = os.Stdout
	} else {
		output = zerolog.ConsoleWriter{
			Out : os.Stdout,
			TimeFormat: zerolog.TimeFieldFormat,
		}
	}

	Logger = zerolog.New(output).
		Level(l).
		With().
		Timestamp().
		Caller().
		Logger()
}
