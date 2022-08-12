package main

import (
	"context"

	"github.com/twothicc/common-go/grpcserver"
	"github.com/twothicc/common-go/logger"
	"github.com/twothicc/datasync/tools/env"
)

var ctx = context.Background()

func main() {
	logger.InitLogger(env.IsTest())

	env.Init(ctx)

	dependencies := initDependencies()

	grpcserver.InitAndRunGrpcServer(ctx, dependencies.ServerConfigs)
}
