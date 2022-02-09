// Package message provides all aspects of the arrebato server for managing messages. Primarily the production and
// consumption of messages.
package message

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/davidsbond/arrebato/internal/command"
	"github.com/davidsbond/arrebato/internal/consumer"
	consumerpb "github.com/davidsbond/arrebato/internal/proto/arrebato/consumer/v1"
	messagecmd "github.com/davidsbond/arrebato/internal/proto/arrebato/message/command/v1"
	messagesvc "github.com/davidsbond/arrebato/internal/proto/arrebato/message/service/v1"
	"github.com/davidsbond/arrebato/internal/proto/arrebato/message/v1"
)

type (
	// The GRPC type is a messagesvc.MessageServiceServer implementation that handles inbound gRPC requests to manage
	// and query Messages.
	GRPC struct {
		executor  Executor
		reader    Reader
		consumers TopicIndexGetter
	}

	// The Executor interface describes types that execute commands related to Message data.
	Executor interface {
		// Execute should perform actions corresponding to the provided command. The returned error should correspond
		// to the issue relevant to the command. For example, a command that creates a new Message should return
		// ErrNoTopic if the topic does not exist.
		Execute(ctx context.Context, cmd command.Command) error
	}

	// The Reader interface describes types that can read messages from a topic starting from a given index.
	Reader interface {
		// Read should start reading messages within a topic starting from a given index, invoking the ReadFunc for
		// each message.
		Read(ctx context.Context, topic string, startIndex uint64, fn ReadFunc) error
	}

	// The TopicIndexGetter interface describes types that can retrieve the current index on a topic for a consumer.
	TopicIndexGetter interface {
		// GetTopicIndex should return the number that describes the current index in the topic that a consumer has
		// read to. It should return ErrNoTopic if the topic does not exist.
		GetTopicIndex(ctx context.Context, topic string, consumerID string) (*consumerpb.TopicIndex, error)
	}
)

// NewGRPC returns a new instance of the GRPC type that will modify Message data via commands sent to the Executor and
// read messages via the Reader implementation.
func NewGRPC(executor Executor, reader Reader, consumers TopicIndexGetter) *GRPC {
	return &GRPC{executor: executor, reader: reader, consumers: consumers}
}

// Register the GRPC service onto the grpc.ServiceRegistrar.
func (svr *GRPC) Register(registrar grpc.ServiceRegistrar) {
	messagesvc.RegisterMessageServiceServer(registrar, svr)
}

// Produce a new message for a topic. Returns a codes.NotFound code if the topic does not exist.
func (svr *GRPC) Produce(ctx context.Context, request *messagesvc.ProduceRequest) (*messagesvc.ProduceResponse, error) {
	cmd := command.New(&messagecmd.CreateMessage{
		Message: &message.Message{
			Topic:     request.GetMessage().GetTopic(),
			Payload:   request.GetMessage().GetPayload(),
			Timestamp: timestamppb.Now(),
		},
	})

	err := svr.executor.Execute(ctx, cmd)
	switch {
	case errors.Is(err, command.ErrNotLeader):
		return nil, status.Error(codes.FailedPrecondition, "this node is not the leader")
	case errors.Is(err, ErrNoTopic):
		return nil, status.Errorf(codes.NotFound, "topic %s does not exist", request.GetMessage().GetTopic())
	case err != nil:
		return nil, status.Errorf(codes.Internal, "failed to create message: %v", err)
	default:
		return &messagesvc.ProduceResponse{}, nil
	}
}

// Consume messages from a topic. Returns a codes.NotFound code if the topic does not exist.
func (svr *GRPC) Consume(request *messagesvc.ConsumeRequest, server messagesvc.MessageService_ConsumeServer) error {
	ctx := server.Context()

	topicIndex, err := svr.consumers.GetTopicIndex(ctx,
		request.GetTopic(),
		request.GetConsumerId(),
	)
	switch {
	case errors.Is(err, consumer.ErrNoTopic):
		return status.Errorf(codes.NotFound, "topic %s does not exist", request.GetTopic())
	case ctx.Err() != nil:
		return status.FromContextError(ctx.Err()).Err()
	case err != nil:
		return status.Errorf(codes.Internal, "failed to get last index for consumer: %v", err)
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return status.FromContextError(ctx.Err()).Err()
		case <-ticker.C:
			// Messages here are read and forwarded to the client from the last index onwards. Once we reach the end
			// of all available messages we wait via the ticker and check again for new messages from the last known
			// index. This allows the client to stay connected and keep getting more messages, without over-polling
			// the data store.
			err := svr.reader.Read(ctx, request.GetTopic(), topicIndex.GetIndex(), func(ctx context.Context, m *message.Message) error {
				resp := &messagesvc.ConsumeResponse{Message: m}

				if err := server.Send(resp); err != nil {
					return fmt.Errorf("failed to send message: %w", err)
				}

				topicIndex.Index = m.GetIndex() + 1
				return nil
			})

			switch {
			case errors.Is(err, ErrNoTopic):
				return status.Errorf(codes.NotFound, "topic %s does not exist", request.GetTopic())
			case errors.Is(err, ErrNoMessages):
				// If no messages have been produced we'll just keep waiting and trying to read.
				continue
			case ctx.Err() != nil:
				return status.FromContextError(ctx.Err()).Err()
			case err != nil:
				return status.Errorf(codes.Internal, "failed to read messages: %v", err)
			}
		}
	}
}
