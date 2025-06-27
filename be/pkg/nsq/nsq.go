package nsq_client

import (
	"context"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/configs"
	"github.com/nsqio/go-nsq"
	"log"
	"time"
)

type (
	Nsq interface {
		Publish(ctx context.Context, topic string, body []byte) error
		Consume(ctx context.Context, topic string) (string, error)
		RegisterConsumer(topic string, handlerFunc func(string)) error
	}

	nsqClient struct {
		pub     *nsq.Producer
		config  *configs.Config
		lookupd string
	}
)

func NewNsqClient(cfg *configs.Config) (*nsqClient, error) {
	nsqConfig := nsq.NewConfig()

	addr := fmt.Sprintf("%s:%s", cfg.Config.Nsq.NSQDHost, cfg.Config.Nsq.NSQDTCPPort)
	producer, err := nsq.NewProducer(addr, nsqConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create NSQ producer: %w", err)
	}

	return &nsqClient{
		pub:     producer,
		config:  cfg,
		lookupd: fmt.Sprintf("%s:%s", cfg.Config.Nsq.NSQDHost, cfg.Config.Nsq.LookupdHttpPort),
	}, nil
}

func (n *nsqClient) Publish(ctx context.Context, topic string, body []byte) error {
	return n.pub.Publish(topic, body)
}

func (n *nsqClient) Consume(ctx context.Context, topic string) (string, error) {
	nsqConfig := nsq.NewConfig()
	log.Println("Creating NSQ consumer for topic:", topic)

	consumer, err := nsq.NewConsumer(topic, "channel", nsqConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create NSQ consumer: %w", err)
	}

	msgChan := make(chan string)
	consumer.AddHandler(nsq.HandlerFunc(func(msg *nsq.Message) error {
		log.Println("message received:", string(msg.Body))
		msgChan <- string(msg.Body)
		return nil
	}))

	lookupAddr := fmt.Sprintf("%s:%s", n.config.Config.Nsq.NSQDHost, n.config.Config.Nsq.LookupdHttpPort)
	log.Println("Connecting to nsqlookupd at", lookupAddr)

	if err := consumer.ConnectToNSQLookupd(lookupAddr); err != nil {
		return "", fmt.Errorf("failed to connect to NSQ lookupd: %w", err)
	}

	log.Println("Connected. Waiting for message...")
	select {
	case msg := <-msgChan:
		log.Println("Got message from channel")
		return msg, nil
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(10 * time.Second):
		return "", fmt.Errorf("timeout waiting for message")
	}
}

func (n *nsqClient) RegisterConsumer(topic string, handlerFunc func(string)) error {
	nsqConfig := nsq.NewConfig()
	log.Println("Creating NSQ consumer for topic:", topic)

	consumer, err := nsq.NewConsumer(topic, "channel", nsqConfig)
	if err != nil {
		return fmt.Errorf("failed to create NSQ consumer: %w", err)
	}

	consumer.AddHandler(nsq.HandlerFunc(func(msg *nsq.Message) error {
		body := string(msg.Body)
		log.Println("NSQ message received:", body)
		handlerFunc(body)
		return nil
	}))

	lookupAddr := fmt.Sprintf("%s:%s", n.config.Config.Nsq.NSQDHost, n.config.Config.Nsq.LookupdHttpPort)
	log.Println("Connecting to nsqlookupd at", lookupAddr)

	if err := consumer.ConnectToNSQLookupd(lookupAddr); err != nil {
		return fmt.Errorf("failed to connect to NSQ lookupd: %w", err)
	}

	log.Println("NSQ consumer registered and running...")
	return nil
}
