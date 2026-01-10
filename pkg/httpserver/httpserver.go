// Package httpserver пакет для настройки и конфигурирования сервера.
package httpserver

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/adexcell/delayed-notifier/pkg/log"
)

type Config struct {
	Addr            string        `mapstructure:"addr" validate:"required"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout" validate:"required,min=1s"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout" validate:"required,min=1s"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout" validate:"required,min=1s"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout" validate:"required,min=1s"`
	MaxHeaderBytes  int           `mapstructure:"max_header_bytes" validate:"required,min=1"`
}

type Server struct {
	server          *http.Server
	log             log.Log
	shutdownTimeout time.Duration
	shutdownCh      chan struct{}
}

func New(handler http.Handler, c Config, l log.Log) *Server {
	return &Server{
		server: &http.Server{
			Addr:           c.Addr,
			Handler:        handler,
			ReadTimeout:    c.ReadTimeout,
			WriteTimeout:   c.WriteTimeout,
			IdleTimeout:    c.IdleTimeout,
			MaxHeaderBytes: c.MaxHeaderBytes,
		},
		log:             l,
		shutdownTimeout: c.ShutdownTimeout,
		shutdownCh:      make(chan struct{}),
	}
}

func (s *Server) Start() error {
	go func() {
		s.log.Info().Str("addr", s.server.Addr).Msg("Server starting")
		err := s.server.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			s.log.Info().Msg("Server stopped gracefully")
		} else if err != nil {
			s.log.Error().Err(err).Msg("server failed")
		}
		close(s.shutdownCh)
	}()
	return nil
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	s.log.Info().Msg("Server shutting down")
	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Error().Err(err).Msg("Server shutdown error")
		return err
	}
	s.log.Info().Msg("Server stopped")
	return nil
}
