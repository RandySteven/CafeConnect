package kafka_client

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/configs"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type (
	Consumer interface {
		ReadMessage(ctx context.Context, topic string, key string) (result string, err error)
	}

	sub struct {
		addr string
		d    *kafka.Dialer
	}
)

func NewConsumer(config *configs.Config) (*sub, error) {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		TLS:       &tls.Config{},
	}

	kafkaConf := config.Config.Kafka

	return &sub{
		addr: fmt.Sprintf("%s:%s", kafkaConf.Host, kafkaConf.Port),
		d:    dialer,
	}, nil
}

func (s *sub) ReadMessage(ctx context.Context, topic string, key string) (string, error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{s.addr},
		Topic:    topic,
		GroupID:  "your-group-id",
		MinBytes: 10e3,
		MaxBytes: 10e6,
		MaxWait:  500 * time.Millisecond,
	})
	defer r.Close()

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			// Context cancellation or read error
			return "", err
		}
		log.Printf("got message: key=%s, value=%s", string(m.Key), string(m.Value))
		if string(m.Key) == key {
			return string(m.Value), nil
		}
		// else continue reading
	}
}
