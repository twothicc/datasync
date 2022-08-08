package grpc

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/twothicc/common-go/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type grpcService struct {
	configs *ServerConfigs
	server  *grpc.Server
}

func InitGrpcService(ctx context.Context, config *ServerConfigs) *grpcService {
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	for _, registerServerHandler := range config.registerServerHandlers {
		registerServerHandler(grpcServer)
	}

	return &grpcService{
		server:  grpcServer,
		configs: config,
	}
}

func (g *grpcService) Run(ctx context.Context) {
	logger.WithContext(ctx).Info("start server")

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", g.configs.port))
	if err != nil {
		logger.WithContext(ctx).Fatal("failed to listen", zap.Error(err))
	}

	if err := g.server.Serve(lis); err != nil {
		logger.WithContext(ctx).Fatal("failed to init grpc server", zap.Error(err))
	}
}

func (g *grpcService) ListenSignals(ctx context.Context) {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-signalChan

	logger.WithContext(ctx).Info("receive signal, stop server", zap.String("signal", sig.String()))
	time.Sleep(1 * time.Second)

	if g.server != nil {
		g.server.GracefulStop()
	}

	logger.Sync()
}
