// Package httpserver пакет для настройки и конфигурирования сервера.
package httpserver

import (
	"net/http"
	"time"
)

type Config struct {
	Port string `mapstructure:"port"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

func New(cfg Config, router http.Handler) *http.Server {
	return &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}
}
