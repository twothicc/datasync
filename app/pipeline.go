package main

import (
	"github.com/twothicc/datasync/handlers/helloworld"
	infra_grpc "github.com/twothicc/datasync/infra/grpc"
	pb "github.com/twothicc/protobuf/datasync/v1"
	"google.golang.org/grpc"
)

type dependencies struct {
	ServerConfigs *infra_grpc.ServerConfigs
}

func initDependencies() *dependencies {
	registerHelloWorldServiceHandler := func(s *grpc.Server) {
		pb.RegisterHelloWorldServiceServer(s, helloworld.NewHelloWorldServer())
	}
	serverConfigs := infra_grpc.GetServerConfigs(
		registerHelloWorldServiceHandler,
	)

	return &dependencies{
		ServerConfigs: serverConfigs,
	}
}
