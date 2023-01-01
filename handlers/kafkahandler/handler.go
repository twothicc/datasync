package kafkahandler

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/twothicc/common-go/logger"
	"github.com/twothicc/datasync/infra/elasticsearch"
	"github.com/twothicc/datasync/tools/shardname"
	"go.uber.org/zap"
	"gopkg.in/Shopify/sarama.v1"
)

type IMessageHandler interface {
	sarama.ConsumerGroupHandler
}

type MessageHandler struct {
	ctx             context.Context
	esClient        *elastic.Client
	bulkHandler     elasticsearch.IBulkHandler
	existingIndexes map[string]bool // assumption that indexes won't be deleted in ES
	mu              *sync.Mutex
}

type Message struct {
	err        error
	Action     string   `json:"action"`
	Table      string   `json:"table"`
	Schema     string   `json:"schema"`
	OldData    []byte   `json:"old_data"`
	NewData    []byte   `json:"new_data"`
	Pk         []string `json:"pk"`
	encoded    []byte
	Ctimestamp uint32 `json:"ctimestamp"`
	Mtimestamp uint32 `json:"mtimestamp"`
}

func NewMessageHandler(ctx context.Context, esClient *elastic.Client) (IMessageHandler, error) {
	if esClient == nil {
		return nil, ErrConstructor.New("[NewMessageHandler]elasticsearch client cannot be nil")
	}

	mh := &MessageHandler{
		ctx:         ctx,
		esClient:    esClient,
		bulkHandler: elasticsearch.InitAndRunBulkHandler(ctx, esClient, BULK_DELAY*time.Second),
	}

	return mh, nil
}

func (mh *MessageHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (mh *MessageHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (mh *MessageHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for rawMsg := range claim.Messages() {
		isError := false

		var msg Message

		if err := json.Unmarshal(rawMsg.Value, &msg); err != nil {
			logger.WithContext(mh.ctx).Error(
				fmt.Sprintf("[MessageHandler.ConsumeClaim]fail to unmarshal message from topic:(%q) partition:(%d) offset:(%d)", rawMsg.Topic, rawMsg.Partition, rawMsg.Offset),
			)

			isError = true
		}

		oldData := make(map[string]interface{})
		if err := json.Unmarshal(msg.OldData, &oldData); err != nil {
			logger.WithContext(mh.ctx).Error("[MessageHandler.ConsumeClaim]fail to unmarshal old data")

			isError = true
		}

		newData := make(map[string]interface{})
		if err := json.Unmarshal(msg.NewData, &newData); err != nil {
			logger.WithContext(mh.ctx).Error("[MessageHandler.ConsumeClaim]fail to unmarshal new data")

			isError = true
		}

		logger.WithContext(mh.ctx).Info(
			"[MessageHandler.ConsumeClaim]received",
			zap.Uint32("ctimestamp", msg.Ctimestamp),
			zap.Uint32("mtimestamp", msg.Mtimestamp),
			zap.String("action", msg.Action),
			zap.String("schema", msg.Schema),
			zap.String("table", msg.Table),
			zap.Strings("pk", msg.Pk),
			zap.Any("old data", oldData),
			zap.Any("new data", newData),
		)

		switch msg.Action {
		case INSERT_ACTION:
			if err := mh.handleInsert(&msg, newData); err != nil {
				if ErrIndex.Is(err) {
					return err
				} else {
					logger.WithContext(mh.ctx).Error(
						"[MessageHandler.ConsumeClaim]fail to insert record",
						zap.Error(err),
					)

					isError = true
				}
			}
		case UPDATE_ACTION:
			if err := mh.handleUpdate(&msg, oldData, newData); err != nil {
				if ErrIndex.Is(err) {
					return err
				} else {
					logger.WithContext(mh.ctx).Error(
						"[MessageHandler.ConsumeClaim]fail to update record",
						zap.Error(err),
					)

					isError = true
				}
			}
		case DELETE_ACTION:
			if err := mh.handleDelete(&msg, oldData); err != nil {
				if ErrIndex.Is(err) {
					return err
				} else {
					logger.WithContext(mh.ctx).Error(
						"[MessageHandler.ConsumeClaim]fail to delete record",
						zap.Error(err),
					)

					isError = true
				}
			}
		default:
			logger.WithContext(mh.ctx).Error(
				"[MessageHandler.ConsumeClaim]unknown action",
				zap.String("action", msg.Action),
			)
		}

		// if no errors in processing, we mark message as done
		// otherwise, it will be processed again until retries is exhausted
		if !isError {
			sess.MarkMessage(rawMsg, "")
		}
	}

	return nil
}

func (mh *MessageHandler) createIndexIfNotExist(indexName string) error {
	mh.mu.Lock()
	defer mh.mu.Unlock()

	if _, ok := mh.existingIndexes[indexName]; !ok {
		existService := elastic.NewIndicesExistsService(mh.esClient).Index([]string{indexName})

		isExist, err := existService.Do(mh.ctx)
		if err != nil {
			logger.WithContext(mh.ctx).Error(
				"[MessageHandler.createIndexIfNotExist]fail to check if index exists",
				zap.String("indexName", indexName),
				zap.Error(err),
			)

			return ErrIndex.New(fmt.Sprintf("[MessageHandler.createIndexIfNotExist]%s", err.Error()))
		}

		if !isExist {
			createService := elastic.NewIndicesCreateService(mh.esClient).Index(indexName)

			createResult, err := createService.Do(mh.ctx)
			if err != nil {
				logger.WithContext(mh.ctx).Error(
					"[MessageHandler.createIndexIfNotExist]fail to create index",
					zap.String("indexName", indexName),
					zap.Error(err),
				)

				return ErrIndex.New(fmt.Sprintf("[MessageHandler.createIndexIfNotExist]%s", err.Error()))
			}

			if !createResult.Acknowledged {
				return ErrIndex.New("[MessageHandler.createIndexIfNotExist]index creation request not acknowledged")
			}

			logger.WithContext(mh.ctx).Info(
				"[MessageHandler.createIndexIfNotExist]successfully created index",
				zap.String("indexName", indexName),
			)
		}

		mh.existingIndexes[indexName] = true
	}

	return nil
}

func (mh *MessageHandler) handleInsert(msg *Message, newData map[string]interface{}) error {
	indexName := shardname.GetShardName(msg.Table, msg.Ctimestamp)
	uniqueId := getUniqueID(msg.Pk)

	if indexErr := mh.createIndexIfNotExist(indexName); indexErr != nil {
		return indexErr
	}

	req := elastic.NewBulkIndexRequest().Index(indexName).Id(uniqueId).Doc(newData)

	mh.bulkHandler.AddBulkRequest(mh.ctx, req)

	logger.WithContext(mh.ctx).Info(
		"[MessageHandler.handleInsert]successfully added insert req",
		zap.String("indexName", indexName),
		zap.String("pk", uniqueId),
	)

	return nil
}

func (mh *MessageHandler) handleDelete(msg *Message, oldData map[string]interface{}) error {
	oldCtimestamp, ok := oldData["ctimestamp"].(uint32)
	if !ok {
		return ErrInvalidCtimestamp.New("[MessageHandler.handleDelete]fail to type assert old record ctimestamp")
	}

	indexName := shardname.GetShardName(msg.Table, oldCtimestamp)
	uniqueId := getUniqueID(msg.Pk)

	if indexErr := mh.createIndexIfNotExist(indexName); indexErr != nil {
		return indexErr
	}

	req := elastic.NewBulkDeleteRequest().Index(indexName).Id(uniqueId)

	mh.bulkHandler.AddBulkRequest(mh.ctx, req)

	logger.WithContext(mh.ctx).Info(
		"[MessageHandler.handleDelete]successfully added delete req",
		zap.String("indexName", indexName),
		zap.String("pk", uniqueId),
	)

	return nil
}

func (mh *MessageHandler) handleUpdate(msg *Message, oldData, newData map[string]interface{}) error {
	oldCtimestamp, ok := oldData["ctimestamp"].(uint32)
	if !ok {
		return ErrInvalidCtimestamp.New("[MessageHandler.handleDelete]fail to type assert old record ctimestamp")
	}

	newCtimestamp, ok := newData["ctimestamp"].(uint32)
	if !ok {
		return ErrInvalidCtimestamp.New("[MessageHandler.handleDelete]fail to type assert new record ctimestamp")
	}

	uniqueId := getUniqueID(msg.Pk)

	if oldCtimestamp == newCtimestamp {
		indexName := shardname.GetShardName(msg.Table, oldCtimestamp)

		if indexErr := mh.createIndexIfNotExist(indexName); indexErr != nil {
			return indexErr
		}

		req := elastic.NewBulkUpdateRequest().Index(indexName).Id(uniqueId).Doc(newData)

		mh.bulkHandler.AddBulkRequest(mh.ctx, req)

		logger.WithContext(mh.ctx).Info(
			"[MessageHandler.handleUpdate]successfully added update req",
			zap.String("indexName", indexName),
			zap.String("pk", uniqueId),
		)
	} else {
		oldIndexName := shardname.GetShardName(msg.Table, oldCtimestamp)
		newIndexName := shardname.GetShardName(msg.Table, newCtimestamp)

		if oldIndexErr := mh.createIndexIfNotExist(oldIndexName); oldIndexErr != nil {
			return oldIndexErr
		}

		if newIndexErr := mh.createIndexIfNotExist(newIndexName); newIndexErr != nil {
			return newIndexErr
		}

		deleteReq := elastic.NewBulkDeleteRequest().Index(oldIndexName).Id(uniqueId)

		mh.bulkHandler.AddBulkRequest(mh.ctx, deleteReq)

		indexReq := elastic.NewBulkIndexRequest().Index(newIndexName).Id(uniqueId).Doc(newData)

		mh.bulkHandler.AddBulkRequest(mh.ctx, indexReq)

		logger.WithContext(mh.ctx).Info(
			"[MessageHandler.handleUpdate]successfully added delete and insert req",
			zap.String("old indexName", oldIndexName),
			zap.String("new indexName", newIndexName),
			zap.String("pk", uniqueId),
		)
	}

	return nil
}
