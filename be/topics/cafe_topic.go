package topics

import (
	"context"

	"github.com/RandySteven/CafeConnect/be/enums"
	topics_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/topics"
	nsq_client "github.com/RandySteven/CafeConnect/be/pkg/nsq"
)

type cafeTopic struct {
	nsq nsq_client.Nsq
}

func (c *cafeTopic) WriteMessage(ctx context.Context, value string) (err error) {
	return c.nsq.Publish(ctx, enums.CafeTopic, []byte(value))
}

func (c *cafeTopic) ReadMessage(ctx context.Context) (value string, err error) {
	return c.nsq.Consume(ctx, enums.CafeTopic)
}

var _ topics_interfaces.CafeTopic = &cafeTopic{}

func newCafeTopic(nsq nsq_client.Nsq) *cafeTopic {
	return &cafeTopic{
		nsq: nsq,
	}
}
