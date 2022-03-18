package arrebato

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	signingsvc "github.com/davidsbond/arrebato/internal/proto/arrebato/signing/service/v1"
)

type (
	// The KeyPair type contains the public and private signing keys for the client to use when producing messages.
	KeyPair struct {
		PublicKey  []byte
		PrivateKey []byte
	}
)

var (
	// ErrSigningKeyExists is the error given when attempting to create a singing key pair for a client that already
	// has them.
	ErrSigningKeyExists = errors.New("signing key exists")

	// ErrNoPublicKey is the error given when querying a public key that does not exist.
	ErrNoPublicKey = errors.New("no public key")
)

// CreateSigningKeyPair attempts to create a new KeyPair for a client.
func (c *Client) CreateSigningKeyPair(ctx context.Context) (KeyPair, error) {
	svc := signingsvc.NewSigningServiceClient(c.cluster.leader())
	resp, err := svc.CreateKeyPair(ctx, &signingsvc.CreateKeyPairRequest{
		ClientId: c.config.ClientID,
	})
	switch {
	case status.Code(err) == codes.AlreadyExists:
		return KeyPair{}, ErrSigningKeyExists
	case status.Code(err) == codes.FailedPrecondition:
		c.cluster.findLeader(ctx)
		return c.CreateSigningKeyPair(ctx)
	case err != nil:
		return KeyPair{}, err
	default:
		return KeyPair{
			PublicKey:  resp.GetKeyPair().GetPublicKey(),
			PrivateKey: resp.GetKeyPair().GetPrivateKey(),
		}, nil
	}
}

// SigningPublicKey attempts to return the public signing key for the specified client identifier.
func (c *Client) SigningPublicKey(ctx context.Context, clientID string) ([]byte, error) {
	svc := signingsvc.NewSigningServiceClient(c.cluster.any())
	resp, err := svc.GetPublicKey(ctx, &signingsvc.GetPublicKeyRequest{
		ClientId: clientID,
	})
	switch {
	case status.Code(err) == codes.NotFound:
		return nil, ErrNoPublicKey
	case err != nil:
		return nil, err
	default:
		return resp.GetPublicKey(), nil
	}
}
