package topics

import (
	topics_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/topics"
	kafka_client "github.com/RandySteven/CafeConnect/be/pkg/kafka"
)

type Topics struct {
	TransactionTopic topics_interfaces.TransactionTopic
	OnboardingTopic  topics_interfaces.OnboardingTopic
}

func NewTopics(pub kafka_client.Publisher, con kafka_client.Consumer) *Topics {
	topic := kafka_client.NewTopic(con, pub)
	return &Topics{
		TransactionTopic: newTransactionTopic(topic),
		OnboardingTopic:  newOnboardingTopic(topic),
	}
}
