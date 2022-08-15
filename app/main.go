package main

import (
	"context"

	"github.com/twothicc/common-go/grpcserver"
	"github.com/twothicc/common-go/logger"
	"github.com/twothicc/datasync/tools/env"
	"go.uber.org/zap/zapcore"
)

var ctx = context.Background()

func main() {
	logger.InitLogger(zapcore.ErrorLevel)

	env.Init(ctx)
	if env.IsTest() {
		logger.InitLogger(zapcore.DebugLevel)
	}

	dependencies := initDependencies()

	grpcserver.InitAndRunGrpcServer(ctx, dependencies.ServerConfigs)
}
