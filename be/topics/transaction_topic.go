package topics

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/enums"
	topics_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/topics"
	kafka_client "github.com/RandySteven/CafeConnect/be/pkg/kafka"
)

type transactionTopic struct {
	topic *kafka_client.Topic
}

func (t *transactionTopic) ReadMessage(ctx context.Context, key string) (result string, err error) {
	return t.topic.Consumer.ReadMessage(ctx, key)
}

func (t *transactionTopic) WriteMessage(ctx context.Context, key string, value string) (err error) {
	return t.topic.Publisher.WriteMessage(ctx, key, value)
}

var _ topics_interfaces.TransactionTopic = &transactionTopic{}

func newTransactionTopic(topic *kafka_client.Topic) *transactionTopic {
	topic.SetTopic(enums.TransactionTopic)
	return &transactionTopic{
		topic: topic,
	}
}
