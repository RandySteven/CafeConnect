package kafka_client

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/configs"
	"github.com/segmentio/kafka-go"
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

func NewConsumer(config *configs.Config) *sub {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		TLS:       &tls.Config{},
	}

	kafkaConf := config.Config.Kafka

	return &sub{
		addr: fmt.Sprintf("%s:%s", kafkaConf.Host, kafkaConf.Port),
		d:    dialer,
	}
}

func (s *sub) ReadMessage(ctx context.Context, topic string, key string) (string, error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{s.addr},
		Topic:     topic,
		Partition: 0,
		MaxBytes:  10e6,
	})
	defer r.Close()

	// Optional: set to start from earliest or latest
	if err := r.SetOffset(kafka.FirstOffset); err != nil {
		return "", err
	}

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			return "", err
		}
		if string(m.Key) == key {
			return string(m.Value), nil
		}
	}
}
