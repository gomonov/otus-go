package internalgrpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/domain"
	"github.com/gomonov/otus-go/hw12_13_14_15_calendar/internal/logger"
	"google.golang.org/grpc"
)

type Server struct {
	server *grpc.Server
	logger Logger
	app    Application
	config GRPCConf
	UnimplementedCalendarServiceServer
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

type GRPCConf struct {
	Host string
	Port string
}

func NewServer(logger *logger.Logger, app Application, config GRPCConf) *Server {
	return &Server{
		logger: logger,
		app:    app,
		config: config,
	}
}

func (s *Server) Start(ctx context.Context) error {
	addr := net.JoinHostPort(s.config.Host, s.config.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.server = grpc.NewServer(
		grpc.UnaryInterceptor(s.loggingInterceptor),
	)
	RegisterCalendarServiceServer(s.server, s)

	go func() {
		s.logger.Info(fmt.Sprintf("gRPC server starting on %s", addr))
		if err := s.server.Serve(lis); err != nil {
			s.logger.Error(fmt.Sprintf("gRPC server failed: %v", err))
		}
	}()

	go func() {
		<-ctx.Done()
		s.Stop(ctx)
	}()

	return nil
}

func (s *Server) Stop(_ context.Context) error {
	s.logger.Info("gRPC server shutting down...")
	if s.server != nil {
		s.server.GracefulStop()
	}
	s.logger.Info("gRPC server stopped")
	return nil
}
