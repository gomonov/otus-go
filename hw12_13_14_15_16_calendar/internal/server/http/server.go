package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/domain"
	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	server *http.Server
	logger Logger
	app    Application
	config ServerConf
}

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

type Application interface {
	CreateEvent(event *domain.Event) error
	GetEvent(id int) (domain.Event, error)
	UpdateEvent(id int, event *domain.Event) error
	DeleteEvent(id int) error
	ListByDay(date time.Time) ([]domain.Event, error)
	ListByWeek(date time.Time) ([]domain.Event, error)
	ListByMonth(date time.Time) ([]domain.Event, error)
}

type ServerConf struct {
	Host string
	Port string
}

func NewServer(logger *logger.Logger, app Application, config ServerConf) *Server {
	return &Server{
		logger: logger,
		app:    app,
		config: config,
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := s.setupRoutes()

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
