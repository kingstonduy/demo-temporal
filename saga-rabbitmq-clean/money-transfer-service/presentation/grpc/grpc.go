package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/lengocson131002/go-clean-core/logger"
	"github.com/ocb/mcs-money-transfer/bootstrap"
	"google.golang.org/grpc"
)

func (s *GrpcServer) GetStartOptions() []GrpcServerStartOption {
	return []GrpcServerStartOption{
		s.WithAccountServer(),
	}
}

func (s *GrpcServer) Start(ctx context.Context) error {
	network, gPort := "tcp", s.port
	lis, err := net.Listen(network, fmt.Sprintf("localhost:%d", gPort))

	if err != nil {
		return err
	}

	opts := s.GetStartOptions()
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return err
		}
	}

	go func() {
		defer func() {
			s.gSrv.GracefulStop()
			s.logger.Info(ctx, "Stop GRPC Server")
		}()
		<-ctx.Done()
	}()

	s.logger.Infof(ctx, "Start GRPC server at port: %v", gPort)
	if err := s.gSrv.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve GRPC %w", err)
	}

	return nil
}

type GrpcServer struct {
	port   int
	logger logger.Logger
	gSrv   *grpc.Server
}

func NewGrpcServer(cfg *bootstrap.ServerConfig, logger logger.Logger) *GrpcServer {
	gSrv := grpc.NewServer()
	return &GrpcServer{
		port:   cfg.GrpcPort,
		logger: logger,
		gSrv:   gSrv,
	}
}
