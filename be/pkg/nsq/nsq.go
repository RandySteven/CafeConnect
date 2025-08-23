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
		RegisterConsumer(topic string, handlerFunc func(context.Context, string)) error
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
	if ctx.Value(topic) != nil {
		log.Println(`context value : `, ctx.Value(topic).(string))
		return ctx.Value(topic).(string), nil
	} else {
		log.Println(`context value : `, nil)
		return "", fmt.Errorf(`failed to consume the topic %s`, topic)
	}
}

func (n *nsqClient) RegisterConsumer(topic string, handlerFunc func(context.Context, string)) error {
	nsqConfig := nsq.NewConfig()
	log.Println("Creating NSQ consumer for topic:", topic)

	consumer, err := nsq.NewConsumer(topic, "channel", nsqConfig)
	if err != nil {
		return fmt.Errorf("failed to create NSQ consumer: %w", err)
	}

	consumer.AddHandler(nsq.HandlerFunc(func(msg *nsq.Message) error {
		body := string(msg.Body)
		ctx := context.WithValue(context.Background(), topic, body)
		ctx, cancel := context.WithTimeout(ctx, time.Second*30)
		defer cancel()

		if err := func() error {
			handlerFunc(ctx, topic)
			return nil
		}(); err != nil {
			log.Println("Error in handlerFunc:", err)
			msg.Requeue(-1)
			return err
		}

		return nil
	}))

	lookupAddr := fmt.Sprintf("%s:%s", n.config.Config.Nsq.NSQDHost, n.config.Config.Nsq.LookupdHttpPort)
	log.Println("Connecting to nsqlookupd at", lookupAddr)

	if err := consumer.ConnectToNSQLookupd(lookupAddr); err != nil {
		return fmt.Errorf("failed to connect to NSQ lookupd: %w", err)
	}

	log.Println("NSQ consumer registered and running... for topic ", topic)
	return nil
}
