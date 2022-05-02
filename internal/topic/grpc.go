// Package topic provides all functionality within arrebato regarding topics. This includes both gRPC, raft and data
// store interactions.
package topic

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	"github.com/davidsbond/arrebato/internal/command"
	nodecmd "github.com/davidsbond/arrebato/internal/proto/arrebato/node/command/v1"
	nodepb "github.com/davidsbond/arrebato/internal/proto/arrebato/node/v1"
	topiccmd "github.com/davidsbond/arrebato/internal/proto/arrebato/topic/command/v1"
	topicsvc "github.com/davidsbond/arrebato/internal/proto/arrebato/topic/service/v1"
	"github.com/davidsbond/arrebato/internal/proto/arrebato/topic/v1"
)

type (
	// The GRPC type is a topicsvc.TopicServiceServer implementation that handles inbound gRPC requests to manage
	// and query Topics.
	GRPC struct {
		executor Executor
		querier  Querier
		nodes    NodeLister
	}

	// The Executor interface describes types that execute commands related to Topic data.
	Executor interface {
		// Execute should perform actions corresponding to the provided command. The returned error should correspond
		// to the issue relevant to the command. For example, a command that creates a new Topic should return
		// ErrTopicExists if the topic already exists.
		Execute(ctx context.Context, cmd command.Command) error
	}

	// The Querier interface describes types that can query Topic data from a data store.
	Querier interface {
		// Get should return the Topic corresponding to the given name. It should return ErrNoTopic if the
		// named topic does not exist. It should return ErrNoTopicInfo if incomplete topic data is found.
		Get(ctx context.Context, name string) (*topic.Topic, error)

		// List should return a slice of all Topics the node has knowledge of. It should return ErrNoTopicInfo if
		// incomplete topic data is found.
		List(ctx context.Context) ([]*topic.Topic, error)
	}

	// The NodeLister interface describes types that can list nodes within the cluster.
	NodeLister interface {
		// List should return all node records in the state.
		List(ctx context.Context) ([]*nodepb.Node, error)
	}
)

// NewGRPC returns a new instance of the GRPC type that will modify Topic data via commands sent to the Executor and
// query Topic data via the Querier implementation. Node data for topic assignment will be obtained via the NodeLister
// implementation.
func NewGRPC(executor Executor, querier Querier, nodes NodeLister) *GRPC {
	return &GRPC{executor: executor, querier: querier, nodes: nodes}
}

// Register the GRPC service onto the grpc.ServiceRegistrar.
func (svr *GRPC) Register(registrar grpc.ServiceRegistrar, health *health.Server) {
	topicsvc.RegisterTopicServiceServer(registrar, svr)
	health.SetServingStatus(topicsvc.TopicService_ServiceDesc.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)
}

// Create a new Topic, assigning it to the node with the least number of assigned topics. Returns a codes.AlreadyExists
// code if the topic already exists.
func (svr *GRPC) Create(ctx context.Context, request *topicsvc.CreateRequest) (*topicsvc.CreateResponse, error) {
	cmd := command.New(&topiccmd.CreateTopic{
		Topic: request.GetTopic(),
	})

	err := svr.executor.Execute(ctx, cmd)
	switch {
	case errors.Is(err, command.ErrNotLeader):
		return nil, status.Error(codes.FailedPrecondition, "this node is not the leader")
	case errors.Is(err, ErrTopicExists):
		return nil, status.Errorf(codes.AlreadyExists, "topic %s already exists", request.GetTopic().GetName())
	case err != nil:
		return nil, status.Errorf(codes.Internal, "failed to create topic %s: %v", request.GetTopic().GetName(), err)
	}

	nodes, err := svr.nodes.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list nodes: %v", err)
	}

	var selected *nodepb.Node
	for _, n := range nodes {
		if len(n.GetTopics()) <= len(selected.GetTopics()) {
			selected = n
		}
	}

	if selected == nil {
		selected = nodes[0]
	}

	cmd = command.New(&nodecmd.AssignTopic{
		TopicName: request.GetTopic().GetName(),
		NodeName:  selected.GetName(),
	})

	err = svr.executor.Execute(ctx, cmd)
	switch {
	case errors.Is(err, command.ErrNotLeader):
		return nil, status.Error(codes.FailedPrecondition, "this node is not the leader")
	case err != nil:
		return nil, status.Errorf(codes.Internal, "failed to assign topic %s: %v", request.GetTopic().GetName(), err)
	}

	return &topicsvc.CreateResponse{}, nil
}

// Delete an existing Topic. Returns a codes.NotFound code if the topic does not exist.
func (svr *GRPC) Delete(ctx context.Context, request *topicsvc.DeleteRequest) (*topicsvc.DeleteResponse, error) {
	cmd := command.New(&topiccmd.DeleteTopic{
		Name: request.GetName(),
	})

	err := svr.executor.Execute(ctx, cmd)
	switch {
	case errors.Is(err, command.ErrNotLeader):
		return nil, status.Error(codes.FailedPrecondition, "this node is not the leader")
	case errors.Is(err, ErrNoTopic):
		return nil, status.Errorf(codes.NotFound, "topic %s does not exist", request.GetName())
	case err != nil:
		return nil, status.Errorf(codes.Internal, "failed to delete topic %s: %v", request.GetName(), err)
	default:
		return &topicsvc.DeleteResponse{}, nil
	}
}

// Get a topic by name. Returns a codes.NotFound code if the topic does not exist, or a codes.InvalidArgument code
// if a topic record is found with incomplete data.
func (svr *GRPC) Get(ctx context.Context, request *topicsvc.GetRequest) (*topicsvc.GetResponse, error) {
	t, err := svr.querier.Get(ctx, request.GetName())
	switch {
	case errors.Is(err, ErrNoTopic):
		return nil, status.Errorf(codes.NotFound, "topic %s does not exist", request.GetName())
	case errors.Is(err, ErrNoTopicInfo):
		return nil, status.Errorf(codes.InvalidArgument, "topic %s is not defined, it should be recreated", request.GetName())
	case err != nil:
		return nil, status.Errorf(codes.Internal, "failed to get topic %s: %v", request.GetName(), err)
	default:
		return &topicsvc.GetResponse{Topic: t}, nil
	}
}

// List all topics known to the node. Returns a codes.FailedPrecondition code if a topic record is found with incomplete
// data.
func (svr *GRPC) List(ctx context.Context, _ *topicsvc.ListRequest) (*topicsvc.ListResponse, error) {
	topics, err := svr.querier.List(ctx)
	switch {
	case errors.Is(err, ErrNoTopicInfo):
		return nil, status.Error(codes.FailedPrecondition, "one or more topics are missing topic info and m ust be recreated")
	case err != nil:
		return nil, status.Errorf(codes.Internal, "failed to list topics: %v", err)
	default:
		return &topicsvc.ListResponse{Topics: topics}, nil
	}
}
