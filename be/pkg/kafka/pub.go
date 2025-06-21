package kafka_client

import (
	"context"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/configs"
	"github.com/segmentio/kafka-go"
	"log"
)

type (
	Publisher interface {
		WriteMessage(ctx context.Context, key string, value string) (err error)
		setTopic(topic string)
	}

	pub struct {
		address string
		w       *kafka.Writer
	}
)

func NewPublisher(config *configs.Config) (*pub, error) {
	kafkaConf := config.Config.Kafka
	addr := fmt.Sprintf("%s:%s", kafkaConf.Host, kafkaConf.Port)

	return &pub{
		address: addr,
	}, nil
}

func (p *pub) setTopic(topic string) {
	w := &kafka.Writer{
		Addr:     kafka.TCP(p.address),
		Balancer: &kafka.LeastBytes{},
		Topic:    topic,
	}
	p.w = w
	log.Println(`success set pub topic `, topic)
}

func (p *pub) WriteMessage(ctx context.Context, key string, value string) (err error) {
	err = p.w.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(key),
			Value: []byte(value),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
