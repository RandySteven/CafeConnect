package kafka_client

import (
	"context"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/configs"
	"github.com/segmentio/kafka-go"
	"log"
	"net"
	"strconv"
)

type (
	Kafka interface {
		RegisterTopics(topics ...string) error
		ReadTopics() []string
		ClearAllTopics() error
	}

	kafkaClient struct {
		conn *kafka.Conn
	}
)

func NewKafkaClient(config *configs.Config) (*kafkaClient, error) {
	kafkaConf := config.Config.Kafka
	address := fmt.Sprintf("%s:%s", kafkaConf.Host, kafkaConf.Port)
	log.Println(address)
	conn, err := kafka.DialContext(context.Background(), kafkaConf.Dial, address)
	if err != nil {
		log.Println(`error during connect `, err)
		return nil, err
	}
	defer conn.Close()
	controller, err := conn.Controller()
	if err != nil {
		log.Println(`error during get controller `, err)
		return nil, err
	}
	connLeader, err := kafka.DialContext(context.Background(), kafkaConf.Dial, net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		log.Println(`error during conn leader `, err)
		return nil, err
	}

	log.Println(`success to connect kafka`)
	return &kafkaClient{
		conn: connLeader,
	}, nil
}

func (k *kafkaClient) RegisterTopics(topics ...string) error {
	topicConfigs := make([]kafka.TopicConfig, len(topics))
	for index, topic := range topics {
		topicConfigs[index] = kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}
	}
	err := k.conn.CreateTopics(topicConfigs...)
	if err != nil {
		return err
	}
	return nil
}

func (k *kafkaClient) ReadTopics() []string {
	partitions, err := k.conn.ReadPartitions()
	if err != nil {
		return nil
	}
	topics := make([]string, len(partitions))

	for index, p := range partitions {
		log.Println(p.Topic)
		topics[index] = p.Topic
	}

	return topics
}

func (k *kafkaClient) ClearAllTopics() error {
	partitions, err := k.conn.ReadPartitions()
	if err != nil {
		return err
	}

	topicSet := make(map[string]struct{})
	for _, p := range partitions {
		topicSet[p.Topic] = struct{}{}
	}

	for topic := range topicSet {
		log.Printf("Please delete topic manually or via CLI/Admin API: %s", topic)
	}

	return nil
}

var _ Kafka = &kafkaClient{}
