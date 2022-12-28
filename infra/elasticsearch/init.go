package elastic

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/twothicc/common-go/logger"
	"github.com/twothicc/datasync/config"
	"github.com/twothicc/datasync/tools/env"
	"go.uber.org/zap"
)

func NewElasticSearchClient(ctx context.Context, elasticCfg *config.ElasticConfig) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: elasticCfg.AddressList,
		Transport: &FastTransport{
			Transport: http.Transport{
				MaxIdleConnsPerHost:   int(elasticCfg.MaxIdleConnsPerHost),
				ResponseHeaderTimeout: time.Duration(elasticCfg.ResponseHeaderTimeout) * time.Millisecond,
				DialContext:           (&net.Dialer{Timeout: time.Nanosecond}).DialContext,
			},
		},
		Username: env.EnvConfigs.ElasticUser,
		Password: env.EnvConfigs.ElasticPass,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		logger.WithContext(ctx).Error("[NewElasticSearchClient]fail to initialize elasticsearch client", zap.Error(err))

		return nil, ErrConstructor.New(fmt.Sprintf("[NewElasticSearchClient]%s", err.Error()))
	}

	logger.WithContext(ctx).Info("[NewElasticSearchClient]successfully initialized elasticsearch client")

	return client, nil
}
