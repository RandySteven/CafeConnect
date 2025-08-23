package topics

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/enums"
	topics_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/topics"
	nsq_client "github.com/RandySteven/CafeConnect/be/pkg/nsq"
)

type onboardingTopic struct {
	nsq nsq_client.Nsq
}

func (o *onboardingTopic) RegisterConsumer(handler func(context.Context, string)) error {
	return o.nsq.RegisterConsumer(enums.OnboardingTopic, handler)
}

func (o *onboardingTopic) WriteMessage(ctx context.Context, value string) (err error) {
	return o.nsq.Publish(ctx, enums.OnboardingTopic, []byte(value))
}

func (o *onboardingTopic) ReadMessage(ctx context.Context) (value string, err error) {
	return o.nsq.Consume(ctx, enums.OnboardingTopic)
}

var _ topics_interfaces.OnboardingTopic = &onboardingTopic{}

func newOnboardingTopic(nsq nsq_client.Nsq) *onboardingTopic {
	return &onboardingTopic{
		nsq: nsq,
	}
}
