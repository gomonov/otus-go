package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	server *http.Server
	logger Logger
	app    Application
	config config.ServerConf
}

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

type Application interface {
}

func NewServer(logger *logger.Logger, app Application, config config.ServerConf) *Server {
	return &Server{
		logger: logger,
		app:    app,
		config: config,
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.helloHandler)
	mux.HandleFunc("/hello", s.helloHandler)

	handler := loggingMiddleware(s.logger, mux)

	s.server = &http.Server{
		Addr:         net.JoinHostPort(s.config.Host, s.config.Port),
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		s.logger.Info(fmt.Sprintf("HTTP server starting on %s", s.server.Addr))

		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error(fmt.Sprintf("HTTP server failed: %v", err))
		}
	}()

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("HTTP server shutting down...")

	if s.server != nil {
		if err := s.server.Shutdown(ctx); err != nil {
			return fmt.Errorf("HTTP server shutdown error: %w", err)
		}
	}

	s.logger.Info("HTTP server stopped")
	return nil
}

func (s *Server) helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/hello" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}
