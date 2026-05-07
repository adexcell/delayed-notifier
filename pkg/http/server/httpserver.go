package httpserver

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// Config holds the configuration parameters for the HTTP server.
type Config struct {
	Host string `default:"localhost" envconfig:"HTTP_HOST"`
	Port string `default:"8080"      envconfig:"HTTP_PORT"`
}

// Server wraps the standard http.Server to provide custom startup and shutdown logic.
type Server struct {
	server *http.Server
}

// New creates and starts a new HTTP server with the provided handler and configuration.
func New(handler http.Handler, c Config) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		Addr:         net.JoinHostPort(c.Host, c.Port),
	}

	s := &Server{
		server: httpServer,
	}

	go s.start()

	log.Info().Msg("http server: started on port: " + httpServer.Addr)

	return s
}

func (s *Server) start() {
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error().Err(err).Msg("http server: ListenAndServe")
	}
}

// Close gracefully shuts down the HTTP server with a timeout.
func (s *Server) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("http server: s.server.Shutdown")
	}

	log.Info().Msg("http server: closed")
}
