package topics

import (
	topics_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/topics"
	nsq_client "github.com/RandySteven/CafeConnect/be/pkg/nsq"
)

type Topics struct {
	TransactionTopic topics_interfaces.TransactionTopic
	OnboardingTopic  topics_interfaces.OnboardingTopic
	MidtransTopic    topics_interfaces.MidtransTopic
	PointTopic       topics_interfaces.PointTopic
}

func NewTopics(nsq nsq_client.Nsq) *Topics {
	return &Topics{
		TransactionTopic: newTransactionTopic(nsq),
		OnboardingTopic:  newOnboardingTopic(nsq),
		MidtransTopic:    newMidtransTopic(nsq),
		PointTopic:       newPointTopic(nsq),
	}
}
