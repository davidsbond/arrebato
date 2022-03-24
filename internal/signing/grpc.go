package signing

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	"github.com/davidsbond/arrebato/internal/clientinfo"
	"github.com/davidsbond/arrebato/internal/command"
	signingcmd "github.com/davidsbond/arrebato/internal/proto/arrebato/signing/command/v1"
	signingsvc "github.com/davidsbond/arrebato/internal/proto/arrebato/signing/service/v1"
	"github.com/davidsbond/arrebato/internal/proto/arrebato/signing/v1"
)

type (
	// The GRPC type is a signingsvc.SigningServiceServer implementation that handles inbound gRPC requests to manage
	// client signing keys.
	GRPC struct {
		publicKeys PublicKeyGetter
		executor   Executor
	}

	// The Executor interface describes types that execute commands related to signing key data.
	Executor interface {
		Execute(ctx context.Context, cmd command.Command) error
	}

	// The PublicKeyGetter interface describes types that can obtain a public key for a client.
	PublicKeyGetter interface {
		Get(ctx context.Context, clientID string) ([]byte, error)
	}
)

// NewGRPC returns a new instance of the GRPC type that will handle inbound gRPC requests for signing key data. Commands
// will be executed via the Executor implementation and public keys will be queried via the PublicKeyGetter implementation.
func NewGRPC(executor Executor, publicKeys PublicKeyGetter) *GRPC {
	return &GRPC{
		publicKeys: publicKeys,
		executor:   executor,
	}
}

// Register the GRPC service onto the grpc.ServiceRegistrar.
func (svr *GRPC) Register(registrar grpc.ServiceRegistrar, healthServer *health.Server) {
	signingsvc.RegisterSigningServiceServer(registrar, svr)
	healthServer.SetServingStatus(signingsvc.SigningService_ServiceDesc.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)
}

// CreateKeyPair attmpts to create a new signing key pair for the client. It returns a codes.FailedPrecondition error code
// if this node is not the leader, or a codes.AlreadyExists error code if the client already has a key pair. The private
// key is not stored by the server and should be kept securely by the client.
func (svr *GRPC) CreateKeyPair(ctx context.Context, _ *signingsvc.CreateKeyPairRequest) (*signingsvc.CreateKeyPairResponse, error) {
	info := clientinfo.FromContext(ctx)

	publicKey, privateKey, err := NewKeyPair()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create key pair: %v", err)
	}

	cmd := command.New(&signingcmd.CreatePublicKey{
		ClientId:  info.ID,
		PublicKey: publicKey,
	})

	err = svr.executor.Execute(ctx, cmd)
	switch {
	case errors.Is(err, command.ErrNotLeader):
		return nil, status.Error(codes.FailedPrecondition, "this node is not the leader")
	case errors.Is(err, ErrPublicKeyExists):
		return nil, status.Errorf(codes.AlreadyExists, "a public key already exists for client %s", info.ID)
	case err != nil:
		return nil, status.Errorf(codes.Internal, "failed to create key pair: %v", err)
	default:
		return &signingsvc.CreateKeyPairResponse{
			KeyPair: &signing.KeyPair{
				PrivateKey: privateKey,
				PublicKey:  publicKey,
			},
		}, nil
	}
}

// GetPublicKey returns the public key for a client. It returns a codes.NotFound error code if there is no public key
// for a client.
func (svr *GRPC) GetPublicKey(ctx context.Context, request *signingsvc.GetPublicKeyRequest) (*signingsvc.GetPublicKeyResponse, error) {
	key, err := svr.publicKeys.Get(ctx, request.GetClientId())
	switch {
	case errors.Is(err, ErrNoPublicKey):
		return nil, status.Errorf(codes.NotFound, "no public key found for client: %s", request.GetClientId())
	case err != nil:
		return nil, status.Errorf(codes.Internal, "failed to get public key: %v", err)
	default:
		return &signingsvc.GetPublicKeyResponse{
			PublicKey: key,
		}, nil
	}
}
