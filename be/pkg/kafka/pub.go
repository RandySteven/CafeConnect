package kafka_client

import (
	"context"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/configs"
	"github.com/segmentio/kafka-go"
)

type (
	Publisher interface {
		WriteMessage(ctx context.Context, topic string, key string, value string) (err error)
	}

	pub struct {
		w *kafka.Writer
	}
)

func NewPublisher(config *configs.Config) *pub {
	kafkaConf := config.Config.Kafka
	w := &kafka.Writer{
		Addr: kafka.TCP(
			fmt.Sprintf("%s:%s", kafkaConf.Host, kafkaConf.Port)),
		Balancer: &kafka.LeastBytes{},
	}

	return &pub{
		w: w,
	}
}

func (p *pub) WriteMessage(ctx context.Context, topic string, key string, value string) (err error) {
	err = p.w.WriteMessages(ctx,
		kafka.Message{
			Topic: topic,
			Key:   []byte(key),
			Value: []byte(value),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
