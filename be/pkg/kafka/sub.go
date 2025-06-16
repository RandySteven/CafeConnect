package kafka_client

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/configs"
	"github.com/RandySteven/CafeConnect/be/enums"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type (
	Consumer interface {
		ReadMessage(ctx context.Context, key string) (result string, err error)
	}

	sub struct {
		addr   string
		d      *kafka.Dialer
		reader *kafka.Reader
	}
)

func NewConsumer(config *configs.Config) (*sub, error) {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		TLS:       &tls.Config{},
	}
	kafkaConf := config.Config.Kafka

	addr := fmt.Sprintf("%s:%s", kafkaConf.Host, kafkaConf.Port)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{addr},
		Topic:    enums.TransactionTopic,
		GroupID:  fmt.Sprintf("test-group-%d", time.Now().UnixNano()),
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	return &sub{
		addr:   addr,
		d:      dialer,
		reader: r,
	}, nil
}

func (s *sub) ReadMessage(ctx context.Context, key string) (string, error) {

	for {
		m, err := s.reader.ReadMessage(ctx)
		if err != nil {
			log.Println(err)
			return "", err
		}
		log.Printf("got message: key=%s, value=%s", string(m.Key), string(m.Value))
		if string(m.Key) == key {
			return string(m.Value), nil
		}
		// else continue reading
	}
}
