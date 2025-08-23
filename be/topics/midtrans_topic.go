package topics

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/enums"
	topics_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/topics"
	nsq_client "github.com/RandySteven/CafeConnect/be/pkg/nsq"
)

type midtransTopic struct {
	nsq nsq_client.Nsq
}

func (m *midtransTopic) RegisterConsumer(handler func(context.Context, string)) error {
	return m.nsq.RegisterConsumer(enums.PaymentMidtransTopic, handler)
}

func (m *midtransTopic) WriteMessage(ctx context.Context, value string) (err error) {
	return m.nsq.Publish(ctx, enums.PaymentMidtransTopic, []byte(value))
}

func (m *midtransTopic) ReadMessage(ctx context.Context) (value string, err error) {
	return m.nsq.Consume(ctx, enums.PaymentMidtransTopic)
}

var _ topics_interfaces.MidtransTopic = &midtransTopic{}

func newMidtransTopic(nsq nsq_client.Nsq) *midtransTopic {
	return &midtransTopic{
		nsq: nsq,
	}
}
