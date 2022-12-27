package main

import (
	"context"

	"github.com/twothicc/common-go/grpcserver"
	"github.com/twothicc/common-go/logger"
	"github.com/twothicc/datasync/config"
	"github.com/twothicc/datasync/handlers/helloworld"
	"github.com/twothicc/datasync/handlers/sync"
	"github.com/twothicc/datasync/tools/env"
	pb "github.com/twothicc/protobuf/datasync/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type dependencies struct {
	ServerConfigs *grpcserver.ServerConfigs
	AppConfigs    *config.Config
}

func initDependencies(ctx context.Context) *dependencies {
	appConfig, err := config.NewConfig("./conf/app.toml")
	if err != nil {
		logger.WithContext(ctx).Error("[initDependencies]fail to load configs", zap.Error(err))
	}

	registerHelloWorldServiceHandler := func(s *grpc.Server) {
		pb.RegisterHelloWorldServiceServer(s, helloworld.NewHelloWorldServer())
	}

	registerSyncServiceHandler := func(s *grpc.Server) {
		pb.RegisterSyncServiceServer(s, sync.NewSyncServer())
	}

	serverConfigs := grpcserver.GetDefaultServerConfigs(
		env.EnvConfigs.ServiceName,
		env.EnvConfigs.Domain,
		env.EnvConfigs.Port,
		env.IsTest(),
		registerHelloWorldServiceHandler,
		registerSyncServiceHandler,
	)

	return &dependencies{
		ServerConfigs: serverConfigs,
		AppConfigs:    appConfig,
	}
}
