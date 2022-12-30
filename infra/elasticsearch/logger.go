package elasticsearch

import (
	"context"

	"github.com/olivere/elastic/v7"
	"github.com/twothicc/common-go/logger"
	"go.uber.org/zap"
)

type elasticErrorLogger struct {
	ctx context.Context
}

type elasticInfoLogger struct {
	ctx context.Context
}

type elasticDebugLogger struct {
	ctx context.Context
}

type iElasticLogger interface {
	elastic.Logger
}

func NewElasticErrorLogger(ctx context.Context) iElasticLogger {
	return &elasticErrorLogger{ctx: ctx}
}

func (eel *elasticErrorLogger) Printf(format string, v ...interface{}) {
	logger.WithContext(eel.ctx).Error("[ElasticClient]", zap.Any("message", v))
}

func NewElasticInfoLogger(ctx context.Context) iElasticLogger {
	return &elasticInfoLogger{ctx: ctx}
}

func (eil *elasticInfoLogger) Printf(format string, v ...interface{}) {
	logger.WithContext(eil.ctx).Info("[ElasticClient]", zap.Any("message", v))
}

func NewElasticDebugLogger(ctx context.Context) iElasticLogger {
	return &elasticDebugLogger{ctx: ctx}
}

func (edl *elasticDebugLogger) Printf(format string, v ...interface{}) {
	logger.WithContext(edl.ctx).Debug("[ElasticClient]", zap.Any("message", v))
}
