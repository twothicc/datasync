package main

import (
	"context"

	"github.com/olivere/elastic/v7"
	"github.com/twothicc/common-go/grpcserver"
	"github.com/twothicc/common-go/logger"
	"github.com/twothicc/datasync/config"
	"github.com/twothicc/datasync/handlers/helloworld"
	"github.com/twothicc/datasync/handlers/synchandler"
	"github.com/twothicc/datasync/infra/elasticsearch"
	"github.com/twothicc/datasync/tools/env"
	pb "github.com/twothicc/protobuf/datasync/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type dependencies struct {
	ServerConfigs *grpcserver.ServerConfigs
	AppConfigs    *config.Config
	EsClient      *elastic.Client
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
		pb.RegisterSyncServiceServer(s, synchandler.NewSyncServer())
	}

	serverConfigs := grpcserver.GetDefaultServerConfigs(
		env.EnvConfigs.ServiceName,
		env.EnvConfigs.Domain,
		env.EnvConfigs.Port,
		env.IsTest(),
		registerHelloWorldServiceHandler,
		registerSyncServiceHandler,
	)

	esClient, err := elasticsearch.NewElasticSearchClient(ctx, &appConfig.ElasticConfig)
	if err != nil {
		logger.WithContext(ctx).Error("[initDependencies]fail to initialize elasticsearch client", zap.Error(err))
	}

	return &dependencies{
		ServerConfigs: serverConfigs,
		AppConfigs:    appConfig,
		EsClient:      esClient,
	}
}
