package arrebato

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	consumersvc "github.com/davidsbond/arrebato/internal/proto/arrebato/consumer/service/v1"
	"github.com/davidsbond/arrebato/internal/proto/arrebato/consumer/v1"
	messagesvc "github.com/davidsbond/arrebato/internal/proto/arrebato/message/service/v1"
	"github.com/davidsbond/arrebato/internal/proto/arrebato/message/v1"
)

type (
	// The Consumer type is used to read messages for a single topic.
	Consumer struct {
		config  ConsumerConfig
		stream  messagesvc.MessageService_ConsumeClient
		cluster *cluster
	}

	// The ConsumerFunc is a function invoked for each message consumed by a Consumer.
	ConsumerFunc func(ctx context.Context, m Message) error

	// The ConsumerConfig type describes configuration values for the Consumer type.
	ConsumerConfig struct {
		Topic      string
		ConsumerID string
	}
)

// NewConsumer returns a new instance of the Consumer type configured to read from a desired topic as a desired
// consumer identifier.
func (c *Client) NewConsumer(ctx context.Context, config ConsumerConfig) (*Consumer, error) {
	conn, err := c.cluster.topicOwner(ctx, config.Topic)
	if err != nil {
		return nil, err
	}

	svc := messagesvc.NewMessageServiceClient(conn)
	stream, err := svc.Consume(ctx, &messagesvc.ConsumeRequest{
		Topic:      config.Topic,
		ConsumerId: config.ConsumerID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start stream: %w", err)
	}

	return &Consumer{
		config:  config,
		stream:  stream,
		cluster: c.cluster,
	}, nil
}

// Consume messages from a topic as a consumer. The last known consumed index is sent to the server on a periodic
// basis so that the consumer can restart from their last known index. Each message consumed will invoke the
// ConsumerFunc. This method blocks until the context is cancelled, the server returns an error or the ConsumerFunc returns
// an error. Close should be called regardless of Consume returning an error.
func (c *Consumer) Consume(ctx context.Context, fn ConsumerFunc) error {
	messages := make(chan *message.Message, 1)
	indexes := make(chan uint64, 1)

	var index uint64
	defer close(messages)
	defer close(indexes)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-messages:
			value, err := msg.GetValue().UnmarshalNew()
			if err != nil {
				return fmt.Errorf("failed to unmarshal value for %s:%v: %w", msg.GetTopic(), msg.GetIndex(), err)
			}

			var key proto.Message
			if msg.GetKey() != nil {
				key, err = msg.GetKey().UnmarshalNew()
				if err != nil {
					return fmt.Errorf("failed to unmarshal key for %s:%v: %w", msg.GetTopic(), msg.GetIndex(), err)
				}
			}

			m := Message{
				Key:   key,
				Value: value,
				Sender: Sender{
					ID:           msg.GetSender().GetId(),
					Verified:     msg.GetSender().GetVerified(),
					KeySignature: msg.GetSender().GetKeySignature(),
				},
			}

			if err = fn(ctx, m); err != nil {
				return err
			}

			indexes <- msg.GetIndex() + 1
		case index = <-indexes:
			req := &consumersvc.SetTopicIndexRequest{
				TopicIndex: &consumer.TopicIndex{
					Topic:      c.config.Topic,
					ConsumerId: c.config.ConsumerID,
					Index:      index,
				},
			}

			svc := consumersvc.NewConsumerServiceClient(c.cluster.leader())
			_, err := svc.SetTopicIndex(ctx, req)
			switch {
			case status.Code(err) == codes.Canceled:
				return context.Canceled
			case status.Code(err) == codes.DeadlineExceeded:
				return context.DeadlineExceeded
			case err != nil:
				return fmt.Errorf("failed to update topic index: %w", err)
			}

		default:
			resp, err := c.stream.Recv()
			switch {
			case status.Code(err) == codes.Canceled:
				return context.Canceled
			case status.Code(err) == codes.DeadlineExceeded:
				return context.DeadlineExceeded
			case err != nil:
				return fmt.Errorf("failed to receive message: %w", err)
			default:
				messages <- resp.GetMessage()
			}
		}
	}
}

// Close the stream of messages. This should be called regardless of Consume returning an error.
func (c *Consumer) Close() error {
	return c.stream.CloseSend()
}
