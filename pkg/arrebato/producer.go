package arrebato

import (
	"context"

	messagesvc "github.com/davidsbond/arrebato/internal/proto/arrebato/message/service/v1"
)

type (
	// The Producer type is responsible for publishing messages onto a single topic.
	Producer struct {
		topic    string
		messages messagesvc.MessageServiceClient
	}
)

// NewProducer returns a new instance of the Producer type that is configured to publish messages for a single
// topic.
func (c *Client) NewProducer(topic string) *Producer {
	return &Producer{
		topic:    topic,
		messages: c.messages,
	}
}

// Produce a message onto the configured topic.
func (p *Producer) Produce(ctx context.Context, m Message) error {
	msg, err := m.toProto()
	if err != nil {
		return err
	}

	msg.Topic = p.topic
	_, err = p.messages.Produce(ctx, &messagesvc.ProduceRequest{
		Message: msg,
	})

	return err
}
