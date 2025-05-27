package apps

import (
	"github.com/RandySteven/CafeConnect/be/configs"
	kafka_client "github.com/RandySteven/CafeConnect/be/pkg/kafka"
	"log"
)

type PubSub struct {
	Writer kafka_client.Publisher
	Reader kafka_client.Consumer
}

func NewPubSub(config *configs.Config) *PubSub {
	w, err := kafka_client.NewPublisher(config)
	if err != nil {
		log.Println(`err create write `, err)
		return nil
	}
	r, err := kafka_client.NewConsumer(config)
	if err != nil {
		log.Println(`err create read `, err)
		return nil
	}
	return &PubSub{
		Writer: w,
		Reader: r,
	}
}
