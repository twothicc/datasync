package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/twothicc/common-go/logger"
	"go.uber.org/zap"
	"gopkg.in/Shopify/sarama.v1"
)

type IMessageHandler interface {
	sarama.ConsumerGroupHandler
}

type SyncMessageHandler struct {
	ctx context.Context
}

func (smh *SyncMessageHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (smh *SyncMessageHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (smh *SyncMessageHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for rawMsg := range claim.Messages() {
		var msg SyncMessage

		if err := json.Unmarshal(rawMsg.Value, &msg); err != nil {
			logger.WithContext(smh.ctx).Error(
				fmt.Sprintf("[SyncMessageHandler.ConsumeClaim]fail to unmarshal message from topic:(%q) partition:(%d) offset:(%d)", rawMsg.Topic, rawMsg.Partition, rawMsg.Offset),
			)

			return ErrUnmarshal.New(fmt.Sprintf("[SyncMessageHandler.ConsumeClaim]%s", err.Error()))
		}

		if err := msg.Handle(smh.ctx); err != nil {
			logger.WithContext(smh.ctx).Error(
				fmt.Sprintf("[SyncMessageHandler.ConsumeClaim]fail to handle message from topic(%s) partition:(%d) offset:(%d)", rawMsg.Topic, rawMsg.Partition, rawMsg.Offset),
				zap.Error(err),
			)

			return ErrConsume.New(fmt.Sprintf("[SyncMessageHandler.ConsumeClaim]%s", err.Error()))
		}

		sess.MarkMessage(rawMsg, "")
	}

	return nil
}

type SyncMessage struct {
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

func (sm *SyncMessage) Handle(ctx context.Context) error {
	oldData := make(map[string]interface{})
	if err := json.Unmarshal(sm.OldData, &oldData); err != nil {
		logger.WithContext(ctx).Error("[SyncServer.Sync]fail to unmarshal old data")
	}

	newData := make(map[string]interface{})
	if err := json.Unmarshal(sm.NewData, &newData); err != nil {
		logger.WithContext(ctx).Error("[SyncServer.Sync]fail to unmarshal new data")
	}

	logger.WithContext(ctx).Info(
		"[SyncMessage.Handle]received",
		zap.Uint32("ctimestamp", sm.Ctimestamp),
		zap.Uint32("mtimestamp", sm.Mtimestamp),
		zap.String("action", sm.Action),
		zap.String("schema", sm.Schema),
		zap.String("table", sm.Table),
		zap.Strings("pk", sm.Pk),
		zap.Any("old data", oldData),
		zap.Any("new data", newData),
	)

	return nil
}
