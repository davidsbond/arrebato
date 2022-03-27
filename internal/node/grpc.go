// Package node provides all functionality within arrebato regarding nodes. This includes the gRPC service implementation.
package node

import (
	"context"

	"github.com/hashicorp/raft"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	nodesvc "github.com/davidsbond/arrebato/internal/proto/arrebato/node/service/v1"
	"github.com/davidsbond/arrebato/internal/proto/arrebato/node/v1"
)

type (
	// The GRPC type is a nodesvc.NodeServiceServer implementation that handles inbound gRPC requests to get
	// information on the node.
	GRPC struct {
		raft Raft
	}

	// The Raft interface describes types that can return information regarding the raft state.
	Raft interface {
		State() raft.RaftState
		GetConfiguration() raft.ConfigurationFuture
	}
)

// NewGRPC returns a new instance of the GRPC type that returns node information based on the Raft state.
func NewGRPC(raft Raft) *GRPC {
	return &GRPC{raft: raft}
}

// Register the GRPC service onto the grpc.ServiceRegistrar.
func (svr *GRPC) Register(registrar grpc.ServiceRegistrar, health *health.Server) {
	nodesvc.RegisterNodeServiceServer(registrar, svr)
	health.SetServingStatus(nodesvc.NodeService_ServiceDesc.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)
}

// Describe the node.
func (svr *GRPC) Describe(_ context.Context, _ *nodesvc.DescribeRequest) (*nodesvc.DescribeResponse, error) {
	future := svr.raft.GetConfiguration()
	if future.Error() != nil {
		return nil, status.Errorf(codes.Internal, "failed to query raft configuration: %v", future.Error())
	}

	config := future.Configuration()
	peers := make([]string, len(config.Servers))
	for i, server := range config.Servers {
		peers[i] = string(server.ID)
	}

	return &nodesvc.DescribeResponse{
		Node: &node.Node{
			Leader: svr.raft.State() == raft.Leader,
			Peers:  peers,
		},
	}, nil
}
