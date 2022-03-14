package arrebato

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	messagesvc "github.com/davidsbond/arrebato/internal/proto/arrebato/message/service/v1"
)

type (
	// The Producer type is responsible for publishing messages onto a single topic.
	Producer struct {
		topic   string
		cluster *cluster
	}
)

// NewProducer returns a new instance of the Producer type that is configured to publish messages for a single
// topic.
func (c *Client) NewProducer(topic string) *Producer {
	return &Producer{
		topic:   topic,
		cluster: c.cluster,
	}
}

// Produce a message onto the configured topic.
func (p *Producer) Produce(ctx context.Context, m Message) error {
	msg, err := m.toProto()
	if err != nil {
		return err
	}

	msg.Topic = p.topic

	svc := messagesvc.NewMessageServiceClient(p.cluster.leader())
	_, err = svc.Produce(ctx, &messagesvc.ProduceRequest{
		Message: msg,
	})
	switch {
	case status.Code(err) == codes.FailedPrecondition:
		p.cluster.findLeader(ctx)
		return p.Produce(ctx, m)
	default:
		return err
	}
}
