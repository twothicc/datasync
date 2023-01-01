package elasticsearch

import (
	"context"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/twothicc/common-go/logger"
	"github.com/twothicc/datasync/tools/env"
	"go.uber.org/zap"
)

type BulkHandler struct {
	bulkReqCh   chan *elastic.BulkableRequest
	bulkRespCh  chan *elastic.BulkResponse
	bulkService *elastic.BulkService
	bulkList    []*elastic.BulkableRequest
	isRunning   bool
}

type IBulkHandler interface {
	AddBulkRequest(ctx context.Context, bulkReq elastic.BulkableRequest)
	Run(ctx context.Context, delay time.Duration)
}

func InitBulkHandler(ctx context.Context, esClient *elastic.Client) IBulkHandler {
	bulkReqCh := make(chan *elastic.BulkableRequest)
	bulkRespCh := make(chan *elastic.BulkResponse)
	bulkList := make([]*elastic.BulkableRequest, 0, MAX_BULK_REQ)
	bulkService := elastic.NewBulkService(esClient)

	return &BulkHandler{
		bulkReqCh:   bulkReqCh,
		bulkRespCh:  bulkRespCh,
		bulkList:    bulkList,
		bulkService: bulkService,
		isRunning:   false,
	}
}

func InitAndRunBulkHandler(ctx context.Context, esClient *elastic.Client, delay time.Duration) IBulkHandler {
	bulkHandler := InitBulkHandler(ctx, esClient)

	bulkHandler.Run(ctx, delay)

	return bulkHandler
}

func (bh *BulkHandler) Run(ctx context.Context, delay time.Duration) {
	if bh.isRunning {
		return
	}

	go func() {
		t := time.NewTicker(delay)
		defer t.Stop()

		for {
			select {
			case req := <-bh.bulkReqCh:
				bh.bulkList = append(bh.bulkList, req)

				if len(bh.bulkList) == cap(bh.bulkList) {
					bulkResp, err := bh.bulkService.Do(ctx)
					if err != nil {
						logger.WithContext(ctx).Error("[BulkHandler.Run]fail to send bulk request", zap.Error(err))
					}

					bh.bulkRespCh <- bulkResp
				}

				bh.bulkService.Add(*req)

				t.Reset(delay)
			case <-t.C:
				t.Reset(delay)
			}
		}
	}()

	go func() {
		for {
			bulkResp := <-bh.bulkRespCh

			if bulkResp.Errors {
				for _, errItem := range bulkResp.Failed() {
					logger.WithContext(ctx).Error("[BulkHandler.Run]fail to handle item", zap.Any("item", errItem))
				}
			}

			if env.IsTest() {
				for _, successItem := range bulkResp.Succeeded() {
					logger.WithContext(ctx).Debug("[BulkHandler.Run]handled item", zap.Any("item", successItem))
				}
			}
		}
	}()

	bh.isRunning = true
}

func (bh *BulkHandler) AddBulkRequest(ctx context.Context, bulkReq elastic.BulkableRequest) {
	bh.bulkReqCh <- &bulkReq
}
