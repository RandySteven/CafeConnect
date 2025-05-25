package consumers

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/enums"
	consumer_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/consumers"
	kafka_client "github.com/RandySteven/CafeConnect/be/pkg/kafka"
	"log"
)

type TransactionConsumer struct {
	consumer kafka_client.Consumer
}

func (t *TransactionConsumer) MidtransTransactionRecord(ctx context.Context) (err error) {
	result, err := t.consumer.ReadMessage(ctx, enums.TransactionTopic, ``)
	if err != nil {
		return err
	}
	log.Println(result)

	return nil
}

var _ consumer_interfaces.TransactionConsumer = &TransactionConsumer{}

func NewTransactionConsumer() *TransactionConsumer {
	return &TransactionConsumer{}
}
