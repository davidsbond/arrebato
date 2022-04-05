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
		PublicKey  []byte `json:"publicKey"`
		PrivateKey []byte `json:"privateKey"`
	}

	// The PublicKey type contains the public key used for verifying signatures for a single client.
	PublicKey struct {
		ClientID  string `json:"clientId"`
		PublicKey []byte `json:"publicKey"`
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
	resp, err := svc.CreateKeyPair(ctx, &signingsvc.CreateKeyPairRequest{})
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
		return resp.GetPublicKey().GetPublicKey(), nil
	}
}

// PublicKeys attempts to return all public keys stored in the server.
func (c *Client) PublicKeys(ctx context.Context) ([]PublicKey, error) {
	svc := signingsvc.NewSigningServiceClient(c.cluster.any())
	resp, err := svc.ListPublicKeys(ctx, &signingsvc.ListPublicKeysRequest{})
	if err != nil {
		return nil, err
	}

	out := make([]PublicKey, len(resp.GetPublicKeys()))
	for i, k := range resp.GetPublicKeys() {
		out[i] = PublicKey{
			ClientID:  k.GetClientId(),
			PublicKey: k.GetPublicKey(),
		}
	}

	return out, nil
}
