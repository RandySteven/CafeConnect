package consumers

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/enums"
	consumer_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/consumers"
	kafka_client "github.com/RandySteven/CafeConnect/be/pkg/kafka"
	"log"
)

type DummyConsumer struct {
	consumer kafka_client.Consumer
}

func (d *DummyConsumer) CheckHealth(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("CheckHealth: context canceled, stopping consumer")
			return
		default:
			result, err := d.consumer.ReadMessage(ctx, enums.DummyTopic, `test-message`)
			if err != nil {
				log.Println(`error reading message:`, err)
				return
			}
			log.Println("Received message:", result)
		}
	}
}

var _ consumer_interfaces.DummyConsumer = &DummyConsumer{}

func newDummyConsumer(
	consumer kafka_client.Consumer,
) *DummyConsumer {
	return &DummyConsumer{
		consumer: consumer,
	}
}
