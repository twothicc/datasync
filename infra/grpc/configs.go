package grpc

import (
	"github.com/twothicc/datasync/tools/env"
	"google.golang.org/grpc"
)

type ServerConfigs struct {
	port                   string
	registerServerHandlers []RegisterServerHandler
}

type RegisterServerHandler func(s *grpc.Server)

func GetServerConfigs(registerServerHandlers ...RegisterServerHandler) *ServerConfigs {
	return &ServerConfigs{
		port:                   env.EnvConfigs.Port,
		registerServerHandlers: registerServerHandlers,
	}
}
