package arrebato

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	messagesvc "github.com/davidsbond/arrebato/internal/proto/arrebato/message/service/v1"
	"github.com/davidsbond/arrebato/internal/proto/arrebato/message/v1"
	"github.com/davidsbond/arrebato/internal/signing"
)

type (
	// The Producer type is responsible for publishing messages onto a single topic.
	Producer struct {
		topic      string
		cluster    *cluster
		signingKey []byte
	}

	// The ProducerConfig type contains configuration values for a Producer.
	ProducerConfig struct {
		// An optional signing key used for messages. When producing messages, if both a message key and this signing
		// key are present, a signature is sent to the server along with the message to verify the message was produced
		// by this client.
		SigningKey []byte

		// The topic to produce messages on.
		Topic string
	}
)

// NewProducer returns a new instance of the Producer type that is configured to publish messages for a single
// topic.
func (c *Client) NewProducer(config ProducerConfig) *Producer {
	return &Producer{
		topic:      config.Topic,
		cluster:    c.cluster,
		signingKey: config.SigningKey,
	}
}

// Produce a message onto the configured topic.
func (p *Producer) Produce(ctx context.Context, m Message) error {
	msg, err := m.toProto()
	if err != nil {
		return err
	}

	msg.Topic = p.topic
	msg.Sender = &message.Sender{}

	// If we have both a message key and a signing key, we'll include the message key signature in the outgoing
	// request metadata. This is used by the server to verify the identity of the client, and tell consumers that
	// the message was indeed produced by this client.
	if m.Key != nil && len(p.signingKey) > 0 {
		signature, err := signing.SignProto(m.Key, p.signingKey)
		if err != nil {
			return fmt.Errorf("failed to sign message key: %w", err)
		}

		msg.Sender.KeySignature = signature
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
