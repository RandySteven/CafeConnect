package consumers

import (
	consumer_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/consumers"
	kafka_client "github.com/RandySteven/CafeConnect/be/pkg/kafka"
)

type Consumers struct {
	DummyConsumer       consumer_interfaces.DummyConsumer
	TransactionConsumer consumer_interfaces.TransactionConsumer
}

func NewConsumers(
	consumer kafka_client.Consumer,
	publisher kafka_client.Publisher,
) *Consumers {
	return &Consumers{
		DummyConsumer:       newDummyConsumer(consumer),
		TransactionConsumer: newTransactionConsumer(consumer),
	}
}
