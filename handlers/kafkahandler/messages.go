package kafkahandler

import (
	"gopkg.in/Shopify/sarama.v1"
)

type IMessageHandler interface {
	sarama.ConsumerGroupHandler
}
