package elasticsearch

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
	"github.com/twothicc/common-go/logger"
	"github.com/twothicc/datasync/config"
	"github.com/twothicc/datasync/tools/env"
	"go.uber.org/zap"
)

func NewElasticSearchClient(ctx context.Context, elasticCfg *config.ElasticConfig) (*elastic.Client, error) {
	client, err := elastic.NewClient(
		newElasticConfig(ctx, elasticCfg)...,
	)
	if err != nil {
		logger.WithContext(ctx).Error("[NewElasticSearchClient]fail to initialize elasticsearch client", zap.Error(err))

		return nil, ErrConstructor.New(fmt.Sprintf("[NewElasticSearchClient]%s", err.Error()))
	}

	logger.WithContext(ctx).Info("[NewElasticSearchClient]successfully initialized elasticsearch client")

	return client, nil
}

func newElasticConfig(ctx context.Context, elasticCfg *config.ElasticConfig) []elastic.ClientOptionFunc {
	config_list := []elastic.ClientOptionFunc{
		elastic.SetURL(elasticCfg.AddressList...),
		elastic.SetBasicAuth(env.EnvConfigs.ElasticUser, env.EnvConfigs.ElasticPass),
		elastic.SetErrorLog(NewElasticErrorLogger(ctx)),
	}

	if env.IsTest() {
		config_list = append(config_list,
			elastic.SetTraceLog(NewElasticDebugLogger(ctx)),
			elastic.SetInfoLog(NewElasticInfoLogger(ctx)),
		)
	}

	return config_list
}
