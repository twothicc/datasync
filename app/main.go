package main

import (
	"context"

	"github.com/twothicc/common-go/logger"
	"github.com/twothicc/datasync/infra/grpc"
	"github.com/twothicc/datasync/tools/env"
)

var ctx = context.Background()

func main() {
	logger.InitLogger(env.IsTest())
	env.Init(ctx)
	grpc.InitAndRunGrpcService(ctx)
}
