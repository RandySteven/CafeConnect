package topics

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/enums"
	topics_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/topics"
	nsq_client "github.com/RandySteven/CafeConnect/be/pkg/nsq"
)

type transactionTopic struct {
	nsq nsq_client.Nsq
}

func (t *transactionTopic) RegisterConsumer(handler func(ctx context.Context, message string)) (err error) {
	return t.nsq.RegisterConsumer(enums.TransactionTopic, handler)
}

func (t *transactionTopic) WriteMessage(ctx context.Context, value string) (err error) {
	return t.nsq.Publish(ctx, enums.TransactionTopic, []byte(value))
}

func (t *transactionTopic) ReadMessage(ctx context.Context) (message string, err error) {
	return t.nsq.Consume(ctx, enums.TransactionTopic)
}

var _ topics_interfaces.TransactionTopic = &transactionTopic{}

func newTransactionTopic(nsq nsq_client.Nsq) *transactionTopic {
	return &transactionTopic{
		nsq: nsq,
	}
}
