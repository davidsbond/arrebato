package acl

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	"github.com/davidsbond/arrebato/internal/command"
	aclcmd "github.com/davidsbond/arrebato/internal/proto/arrebato/acl/command/v1"
	aclsvc "github.com/davidsbond/arrebato/internal/proto/arrebato/acl/service/v1"
	"github.com/davidsbond/arrebato/internal/proto/arrebato/acl/v1"
)

type (
	// The GRPC type is an aclsvc.ACLServiceServer implementation that handles inbound gRPC requests to manage
	// and query the server ACL.
	GRPC struct {
		acl      Getter
		executor Executor
	}

	// The Getter interface describes types that can return the current server ACL.
	Getter interface {
		Get(ctx context.Context) (*acl.ACL, error)
	}

	// The Executor interface describes types that execute commands related to ACL data.
	Executor interface {
		Execute(ctx context.Context, cmd command.Command) error
	}
)

// NewGRPC returns a new instance of the GRPC type that will modify ACL data via commands sent to the Executor and
// query ACL data via the Getter implementation.
func NewGRPC(executor Executor, acl Getter) *GRPC {
	return &GRPC{
		acl:      acl,
		executor: executor,
	}
}

// Register the GRPC service onto the grpc.ServiceRegistrar.
func (svr *GRPC) Register(registrar grpc.ServiceRegistrar, health *health.Server) {
	aclsvc.RegisterACLServiceServer(registrar, svr)
	health.SetServingStatus(aclsvc.ACLService_ServiceDesc.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)
}

// Set the server's new ACL state.
func (svr *GRPC) Set(ctx context.Context, request *aclsvc.SetRequest) (*aclsvc.SetResponse, error) {
	normalized, err := Normalize(ctx, request.GetAcl())
	switch {
	case errors.Is(err, ErrInvalidACL):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case err != nil:
		return nil, status.Errorf(codes.Internal, "failed to normalize ACL: %v", err)
	}

	cmd := command.New(&aclcmd.SetACL{
		Acl: normalized,
	})
	err = svr.executor.Execute(ctx, cmd)
	switch {
	case errors.Is(err, command.ErrNotLeader):
		return nil, status.Error(codes.FailedPrecondition, "this node is not the leader")
	case err != nil:
		return nil, status.Errorf(codes.Internal, "failed to set ACL: %v", err)
	default:
		return &aclsvc.SetResponse{}, nil
	}
}

// Get the server's current ACL state.
func (svr *GRPC) Get(ctx context.Context, _ *aclsvc.GetRequest) (*aclsvc.GetResponse, error) {
	result, err := svr.acl.Get(ctx)
	switch {
	case errors.Is(err, ErrNoACL):
		return nil, status.Error(codes.NotFound, "no acl has been configured on the server")
	case err != nil:
		return nil, status.Errorf(codes.Internal, "failed to query ACL: %v", err)
	}

	normalized, err := Normalize(ctx, result)
	switch {
	case errors.Is(err, ErrInvalidACL):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case err != nil:
		return nil, status.Errorf(codes.Internal, "failed to normalize ACL: %v", err)
	default:
		return &aclsvc.GetResponse{
			Acl: normalized,
		}, nil
	}
}
