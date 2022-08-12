package main

import (
	"github.com/twothicc/common-go/grpcserver"
	"github.com/twothicc/datasync/handlers/helloworld"
	"github.com/twothicc/datasync/tools/env"
	pb "github.com/twothicc/protobuf/datasync/v1"
	"google.golang.org/grpc"
)

type dependencies struct {
	ServerConfigs *grpcserver.ServerConfigs
}

func initDependencies() *dependencies {
	registerHelloWorldServiceHandler := func(s *grpc.Server) {
		pb.RegisterHelloWorldServiceServer(s, helloworld.NewHelloWorldServer())
	}
	serverConfigs := grpcserver.GetDefaultServerConfigs(
		env.EnvConfigs.ServiceName,
		env.EnvConfigs.Domain,
		env.EnvConfigs.Port,
		env.IsTest(),
		registerHelloWorldServiceHandler,
	)

	return &dependencies{
		ServerConfigs: serverConfigs,
	}
}
