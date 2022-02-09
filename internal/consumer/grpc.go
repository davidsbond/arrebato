// Package consumer provides all aspects of the arrebato server for managing consumers. Primarily the storage of consumer
// indexes. Consumer indexes indicate the location within a topic for a given consumer identifier.
package consumer

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/davidsbond/arrebato/internal/command"
	consumercmd "github.com/davidsbond/arrebato/internal/proto/arrebato/consumer/command/v1"
	consumersvc "github.com/davidsbond/arrebato/internal/proto/arrebato/consumer/service/v1"
	consumerpb "github.com/davidsbond/arrebato/internal/proto/arrebato/consumer/v1"
)

type (
	// The GRPC type is a consumersvc.ConsumerServiceServer implementation that handles inbound gRPC requests to manage
	// consumer state.
	GRPC struct {
		executor Executor
	}

	// The Executor interface describes types that execute commands related to Topic data.
	Executor interface {
		// Execute should perform actions corresponding to the provided command. The returned error should correspond
		// to the issue relevant to the command.
		Execute(ctx context.Context, cmd command.Command) error
	}
)

// NewGRPC returns a new instance of the GRPC type that will modify consumer data via commands sent to the Executor.
func NewGRPC(executor Executor) *GRPC {
	return &GRPC{executor: executor}
}

// Register the GRPC service onto the grpc.ServiceRegistrar.
func (svr *GRPC) Register(registrar grpc.ServiceRegistrar) {
	consumersvc.RegisterConsumerServiceServer(registrar, svr)
}

// SetTopicIndex handles an inbound gRPC request to set the current position of a consumer on a topic. Returns a
// codes.FailedPrecondition error code if this server is not the leader, or a codes.NotFound error code if the presented
// topic does not exist.
func (svr *GRPC) SetTopicIndex(ctx context.Context, request *consumersvc.SetTopicIndexRequest) (*consumersvc.SetTopicIndexResponse, error) {
	cmd := command.New(&consumercmd.SetTopicIndex{
		TopicIndex: &consumerpb.TopicIndex{
			Topic:      request.GetTopicIndex().GetTopic(),
			ConsumerId: request.GetTopicIndex().GetConsumerId(),
			Index:      request.GetTopicIndex().GetIndex(),
			Timestamp:  timestamppb.Now(),
		},
	})

	err := svr.executor.Execute(ctx, cmd)
	switch {
	case errors.Is(err, command.ErrNotLeader):
		return nil, status.Error(codes.FailedPrecondition, "this node is not the leader")
	case errors.Is(err, ErrNoTopic):
		return nil, status.Errorf(codes.NotFound, "topic %s does not exist", request.GetTopicIndex().GetTopic())
	case err != nil:
		return nil, status.Errorf(codes.Internal, "failed to set topic index: %v", err)
	default:
		return &consumersvc.SetTopicIndexResponse{}, nil
	}
}
