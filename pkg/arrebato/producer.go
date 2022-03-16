package arrebato

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	messagesvc "github.com/davidsbond/arrebato/internal/proto/arrebato/message/service/v1"
	"github.com/davidsbond/arrebato/internal/signing"
)

type (
	// The Producer type is responsible for publishing messages onto a single topic.
	Producer struct {
		topic      string
		cluster    *cluster
		privateKey []byte
	}
)

// NewProducer returns a new instance of the Producer type that is configured to publish messages for a single
// topic.
func (c *Client) NewProducer(topic string) *Producer {
	return &Producer{
		topic:      topic,
		cluster:    c.cluster,
		privateKey: c.config.MessageSigningKey,
	}
}

// Produce a message onto the configured topic.
func (p *Producer) Produce(ctx context.Context, m Message) error {
	msg, err := m.toProto()
	if err != nil {
		return err
	}

	msg.Topic = p.topic

	// If we have both a message key and a private key, we'll include the message key signature in the outgoing
	// request metadata. This is used by the server to verify the identity of the client, and tell consumers that
	// the message was indeed produced by this client.
	if m.Key != nil && len(p.privateKey) > 0 {
		signature, err := signing.SignProto(m.Key, p.privateKey)
		if err != nil {
			return fmt.Errorf("failed to sign message key: %w", err)
		}

		ctx = metadata.AppendToOutgoingContext(ctx, "X-Key-Signature", string(signature))
	}

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
