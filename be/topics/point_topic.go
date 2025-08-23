package topics

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/enums"
	topics_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/topics"
	nsq_client "github.com/RandySteven/CafeConnect/be/pkg/nsq"
)

type pointTopic struct {
	nsq nsq_client.Nsq
}

func (p *pointTopic) WriteMessage(ctx context.Context, value string) (err error) {
	return p.nsq.Publish(ctx, enums.UserPointTopic, []byte(value))
}

func (p *pointTopic) ReadMessage(ctx context.Context) (value string, err error) {
	return p.nsq.Consume(ctx, enums.UserPointTopic)
}

var _ topics_interfaces.PointTopic = &pointTopic{}

func newPointTopic(nsq nsq_client.Nsq) *pointTopic {
	return &pointTopic{
		nsq: nsq,
	}
}
