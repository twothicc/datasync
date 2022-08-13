package helloworld

import (
	"context"

	"github.com/twothicc/common-go/logger"
	pb "github.com/twothicc/protobuf/datasync/v1"
)

type helloWorldServer struct {
	pb.UnimplementedHelloWorldServiceServer
}

func NewHelloWorldServer() *helloWorldServer {
	return &helloWorldServer{}
}

func (h *helloWorldServer) HelloWorld(
	ctx context.Context,
	req *pb.HelloWorldRequest,
) (*pb.HelloWorldResponse, error) {
	logger.WithContext(ctx).Info("Hello World Request Received")

	return &pb.HelloWorldResponse{
		Msg: "hello world",
	}, nil
}
