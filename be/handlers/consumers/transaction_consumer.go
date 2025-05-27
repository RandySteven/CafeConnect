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
	//1. get midtrans request
	result, err := t.consumer.ReadMessage(ctx, enums.TransactionTopic, ``)
	if err != nil {
		return err
	}
	log.Println(result)
	//2. hit midtrans api

	//3. if midtrans == SUCCESS
	// 3.1. update transaction status from PENDING to SUCCESS
	// 3.2 else try retry mechanism TODO

	return nil
}

var _ consumer_interfaces.TransactionConsumer = &TransactionConsumer{}

func newTransactionConsumer(consumer kafka_client.Consumer) *TransactionConsumer {
	return &TransactionConsumer{
		consumer: consumer,
	}
}
