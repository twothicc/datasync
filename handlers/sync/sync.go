package sync

import (
	"context"
	"encoding/json"

	"github.com/twothicc/common-go/logger"
	pb "github.com/twothicc/protobuf/datasync/v1"
	"go.uber.org/zap"
)

type syncServer struct {
	pb.UnimplementedSyncServiceServer
}

func NewSyncServer() *syncServer {
	return &syncServer{}
}

func (s *syncServer) Sync(
	ctx context.Context,
	req *pb.SyncRequest,
) (*pb.SyncResponse, error) {
	oldData := make(map[string]interface{})
	if err := json.Unmarshal(req.OldData, &oldData); err != nil {
		logger.WithContext(ctx).Error("[SyncServer.Sync]fail to unmarshal old data")
	}

	newData := make(map[string]interface{})
	if err := json.Unmarshal(req.NewData, &newData); err != nil {
		logger.WithContext(ctx).Error("[SyncServer.Sync]fail to unmarshal new data")
	}

	logger.WithContext(ctx).Debug(
		"[SyncServer.Sync]received",
		zap.Uint32("ctimestamp", req.Ctimestamp),
		zap.Uint32("mtimestamp", req.Mtimestamp),
		zap.String("action", req.Action),
		zap.String("schema", req.Schema),
		zap.String("table", req.Table),
		zap.Strings("pk", req.Pk),
		zap.Any("old data", oldData),
		zap.Any("new data", newData),
	)

	return &pb.SyncResponse{
		Msg: "successfully received",
	}, nil
}
