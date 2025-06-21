package kafka_client

type Topic struct {
	Publisher
	Consumer
}

func NewTopic(consumer Consumer, publisher Publisher) *Topic {
	return &Topic{
		Publisher: publisher,
		Consumer:  consumer,
	}
}

func (t *Topic) SetTopic(topic string) {
	t.Consumer.setTopic(topic)
	t.Publisher.setTopic(topic)
}
