package topics

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/enums"
	topics_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/topics"
	kafka_client "github.com/RandySteven/CafeConnect/be/pkg/kafka"
)

type onboardingTopic struct {
	topic *kafka_client.Topic
}

func (o *onboardingTopic) ReadMessage(ctx context.Context, key string) (result string, err error) {
	return o.topic.Consumer.ReadMessage(ctx, key)
}

func (o *onboardingTopic) WriteMessage(ctx context.Context, key string, value string) (err error) {
	return o.topic.Publisher.WriteMessage(ctx, key, value)
}

var _ topics_interfaces.OnboardingTopic = &onboardingTopic{}

func newOnboardingTopic(topic *kafka_client.Topic) *onboardingTopic {
	topic.SetTopic(enums.OnboardingTopic)
	return &onboardingTopic{
		topic: topic,
	}
}
