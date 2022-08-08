package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/twothicc/common-go/logger"
	"github.com/twothicc/datasync/handlers/helloworld"
	"github.com/twothicc/datasync/tools/env"
	pb "github.com/twothicc/protobuf/datasync/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitAndRunGrpcService(ctx context.Context) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", env.EnvConfigs.Port))
	if err != nil {
		logger.WithContext(ctx).Fatal("failed to listen", zap.Error(err))
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	registerServers(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		logger.WithContext(ctx).Fatal("failed to init grpc server", zap.Error(err))
		panic("fail to init grpc server")
	}
}

func registerServers(server *grpc.Server) {
	pb.RegisterHelloWorldServiceServer(server, helloworld.NewHelloWorldServer())
}
