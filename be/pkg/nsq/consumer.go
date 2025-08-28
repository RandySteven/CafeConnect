package nsq_client

import (
	"fmt"
	"log"
	"context"
)

func (n *nsqClient) Consume(ctx context.Context, topic string) (string, error) {
	if ctx.Value(topic) != nil {
		log.Println(`context value : `, ctx.Value(topic).(string))
		return ctx.Value(topic).(string), nil
	} else {
		log.Println(`context value : `, nil)
		return "", fmt.Errorf(`failed to consume the topic %s`, topic)
	}
}
