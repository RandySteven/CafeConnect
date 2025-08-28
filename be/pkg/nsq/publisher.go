package nsq_client

import "context"

func (n *nsqClient) Publish(ctx context.Context, topic string, body []byte) error {
	return n.pub.Publish(topic, body)
}
