package kafka

import (
	"context"
	"fmt"
	"strings"

	"github.com/twothicc/common-go/logger"
	"github.com/twothicc/datasync/config"
	"go.uber.org/zap"
	sarama "gopkg.in/Shopify/sarama.v1"
)

type MessageConsumer struct {
	consumerGroup sarama.ConsumerGroup
	Close         CloseMessageConsumer
}

type CloseMessageConsumer func() error

type IMessageConsumer interface {
	ConsumeTopics(ctx context.Context, topics []string)
}

// MessageConsumer initialized by this constructor must be closed by calling its Close() method
// This is to prevent memory leakage
func NewMessageConsumer(ctx context.Context, kafkaCfg *config.KafkaConfig) (IMessageConsumer, error) {
	version, err := sarama.ParseKafkaVersion(kafkaCfg.Version)
	if err != nil {
		logger.WithContext(ctx).Error("[NewMessageConsumer]fail to parse kafka version", zap.Error(err))

		return nil, ErrConstructor.New(fmt.Sprintf("[NewMessageConsumer]%s", err.Error()))
	}

	saramaCfg := sarama.NewConfig()

	saramaCfg.Consumer.Return.Errors = true
	saramaCfg.Version = version

	switch kafkaCfg.Assignor {
	case "roundrobin":
		saramaCfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	case "range":
		saramaCfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	default:
		logger.WithContext(ctx).Error("[NewMessageConsumer]invalid consumer group rebalance strategy")

		return nil, ErrConstructor.New(fmt.Sprintf("[NewMessageConsumer]%s", err.Error()))
	}

	client, err := sarama.NewClient(kafkaCfg.BrokerList, saramaCfg)
	if err != nil {
		logger.WithContext(ctx).Error("[NewMessageConsumer]fail to initialize kafka client")

		return nil, ErrConstructor.New(fmt.Sprintf("[NewMessageConsumer]%s", err.Error()))
	}

	consumerGroup, err := sarama.NewConsumerGroupFromClient(kafkaCfg.ConsumerGroup, client)
	if err != nil {
		logger.WithContext(ctx).Error(
			fmt.Sprintf("[NewMessageConsumer]fail to initialize consumer group(%s)", kafkaCfg.ConsumerGroup),
		)

		return nil, ErrConstructor.New(fmt.Sprintf("[NewMessageConsumer]%s", err.Error()))
	}

	return &MessageConsumer{
		consumerGroup: consumerGroup,
		Close: func() error {
			if err := consumerGroup.Close(); err != nil {
				logger.WithContext(ctx).Error("[MessageConsumer.Close]fail to close consumer group", zap.Error(err))

				return ErrConsume.New(fmt.Sprintf("[MessageConsumer.Close]%s", err.Error()))
			}

			if err := client.Close(); err != nil {
				logger.WithContext(ctx).Error("[MessageConsumer.Close]fail to close kafka client", zap.Error(err))

				return ErrConsume.New(fmt.Sprintf("[MessageConsumer.Close]%s", err.Error()))
			}

			logger.WithContext(ctx).Info("[MessageConsumer.Close]successfully closed message consumer")

			return nil
		},
	}, nil
}

// Close() is not required because once ConsumeTopics is invoked, MessageConsumer will consume forever
func (mc *MessageConsumer) ConsumeTopics(ctx context.Context, topics []string) {
	go func() {
		for {
			handler := &SyncMessageHandler{
				ctx: ctx,
			}

			if err := mc.consumerGroup.Consume(ctx, topics, handler); err != nil {
				logger.WithContext(ctx).Error(
					fmt.Sprintf(
						"[MessageConsumer.ConsumeTopics]fail to consume topics[%s]",
						strings.Join(topics, ", "),
					),
				)
			}
		}
	}()
}
