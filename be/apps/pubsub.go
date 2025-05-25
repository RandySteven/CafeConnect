package apps

import (
	"github.com/RandySteven/CafeConnect/be/configs"
	kafka_client "github.com/RandySteven/CafeConnect/be/pkg/kafka"
)

type PubSub struct {
	Writer kafka_client.Publisher
	Reader kafka_client.Consumer
}

func NewPubSub(config *configs.Config) *PubSub {
	w := kafka_client.NewPublisher(config)
	r := kafka_client.NewConsumer(config)
	return &PubSub{
		Writer: w,
		Reader: r,
	}
}
