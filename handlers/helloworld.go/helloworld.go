package helloworld

import (
	"context"

	"github.com/google/uuid"
	"github.com/twothicc/common-go/logger"
	pb "github.com/twothicc/protobuf/datasync/v1"
	"go.uber.org/zap"
)

type helloWorldServer struct {
	pb.UnimplementedHelloWorldServiceServer
}

func NewHelloWorldServer() *helloWorldServer {
	return &helloWorldServer{}
}

func (h *helloWorldServer) HelloWorld(
	ctx context.Context,
) (*pb.HelloWorldResponse, error) {

	reqId, _ := uuid.NewRandom()
	ctx = logger.NewLogContext(ctx, zap.Stringer("reqId", reqId))

	logger.WithContext(ctx).Info("Hello World Request Received")

	return &pb.HelloWorldResponse{
		Msg: "hello world",
	}, nil
}
