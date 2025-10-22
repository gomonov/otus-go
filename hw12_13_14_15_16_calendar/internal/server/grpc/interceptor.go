package internalgrpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

func (s *Server) loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)

	s.logger.Info(fmt.Sprintf("gRPC method: %s, duration: %v", info.FullMethod, duration))

	return resp, err
}
