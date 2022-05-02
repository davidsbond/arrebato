// Package node provides all functionality within arrebato regarding nodes. This includes the gRPC service implementation.
package node

import (
	"context"
	"io"
	"time"

	"github.com/hashicorp/raft"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	nodesvc "github.com/davidsbond/arrebato/internal/proto/arrebato/node/service/v1"
	"github.com/davidsbond/arrebato/internal/proto/arrebato/node/v1"
)

type (
	// The GRPC type is a nodesvc.NodeServiceServer implementation that handles inbound gRPC requests to get
	// information on the node.
	GRPC struct {
		raft    Raft
		nodes   Lister
		name    string
		backup  io.WriterTo
		version string
	}

	// The Raft interface describes types that can return information regarding the raft state.
	Raft interface {
		State() raft.RaftState
	}

	// The Lister interface describes types that can list nodes within the state.
	Lister interface {
		// List should return all nodes within the cluster, including the current one.
		List(ctx context.Context) ([]*node.Node, error)
	}

	// The Info type contains node information served from the gRPC service.
	Info struct {
		Name    string
		Version string
	}
)

// NewGRPC returns a new instance of the GRPC type that returns node information based on the Raft state.
func NewGRPC(raft Raft, info Info, backup io.WriterTo, nodes Lister) *GRPC {
	return &GRPC{
		raft:    raft,
		name:    info.Name,
		version: info.Version,
		backup:  backup,
		nodes:   nodes,
	}
}

// Register the GRPC service onto the grpc.ServiceRegistrar.
func (svr *GRPC) Register(registrar grpc.ServiceRegistrar, health *health.Server) {
	nodesvc.RegisterNodeServiceServer(registrar, svr)
	health.SetServingStatus(nodesvc.NodeService_ServiceDesc.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)
}

// Describe the node.
func (svr *GRPC) Describe(ctx context.Context, _ *nodesvc.DescribeRequest) (*nodesvc.DescribeResponse, error) {
	nodes, err := svr.nodes.List(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list nodes")
	}

	peers := make([]string, 0)
	topics := make([]string, 0)
	for _, n := range nodes {
		if n.GetName() == svr.name {
			topics = n.GetTopics()
			continue
		}

		peers = append(peers, n.GetName())
	}

	return &nodesvc.DescribeResponse{
		Node: &node.Node{
			Name:    svr.name,
			Topics:  topics,
			Leader:  svr.raft.State() == raft.Leader,
			Peers:   peers,
			Version: svr.version,
		},
	}, nil
}

// Watch the node state, writing data to the server stream when the leadership or known peers changes.
func (svr *GRPC) Watch(_ *nodesvc.WatchRequest, server nodesvc.NodeService_WatchServer) error {
	ctx := server.Context()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var lastState *node.Node
	for {
		select {
		case <-ctx.Done():
			return status.FromContextError(ctx.Err()).Err()
		case <-ticker.C:
			resp, err := svr.Describe(ctx, &nodesvc.DescribeRequest{})
			if err != nil {
				return err
			}

			state := resp.GetNode()
			if proto.Equal(lastState, state) {
				continue
			}

			if err = server.Send(&nodesvc.WatchResponse{Node: state}); err != nil {
				return err
			}

			lastState = state
		}
	}
}

type (
	backupWriter struct {
		stream nodesvc.NodeService_BackupServer
	}
)

func (bw *backupWriter) Write(p []byte) (int, error) {
	err := bw.stream.Send(&nodesvc.BackupResponse{
		Data: p,
	})

	return len(p), err
}

// Backup the server state, writing it to the outbound stream.
func (svr *GRPC) Backup(_ *nodesvc.BackupRequest, stream nodesvc.NodeService_BackupServer) error {
	if _, err := svr.backup.WriteTo(&backupWriter{stream: stream}); err != nil {
		return status.Errorf(codes.Internal, "failed to perform backup: %v", err)
	}

	return nil
}
