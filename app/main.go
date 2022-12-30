package main

import (
	"context"

	"github.com/twothicc/common-go/grpcserver"
	"github.com/twothicc/common-go/logger"
	"github.com/twothicc/datasync/handlers/kafkahandler/sync"
	"github.com/twothicc/datasync/infra/kafka"
	"github.com/twothicc/datasync/tools/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ctx = context.Background()

func main() {
	logger.InitLogger(zapcore.ErrorLevel)

	env.Init(ctx)

	if env.IsTest() {
		logger.InitLogger(zapcore.DebugLevel)
	}

	dependencies := initDependencies(ctx)

	consumer, err := kafka.NewMessageConsumer(ctx, &dependencies.AppConfigs.KafkaConfig)
	if err != nil {
		logger.WithContext(ctx).Error("[Main]fail to initialize message consumer", zap.Error(err))

		return
	}

	syncMessageHandler, err := sync.NewSyncMessageHandler(ctx, dependencies.EsClient)
	if err != nil {
		logger.WithContext(ctx).Error("[Main]fail to initialize sync message handler", zap.Error(err))

		return
	}

	consumer.ConsumeTopics(ctx, dependencies.AppConfigs.KafkaConfig.Topics, syncMessageHandler)

	grpcserver.InitAndRunGrpcServer(ctx, dependencies.ServerConfigs)
}
