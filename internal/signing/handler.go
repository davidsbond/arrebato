package signing

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"

	signingcmd "github.com/davidsbond/arrebato/internal/proto/arrebato/signing/command/v1"
)

type (
	// The Handler type is responsible for handling commands sent to the server regarding signing keys.
	Handler struct {
		publicKeys PublicKeyCreator
		logger     hclog.Logger
	}

	// The PublicKeyCreator interface describes types that can create public key records for a client.
	PublicKeyCreator interface {
		// Create a public key for a client, should return ErrPublicKeyExists if the client already has a public
		// key.
		Create(ctx context.Context, clientID string, publicKey []byte) error
	}
)

// NewHandler returns a new instance of the Handler type that will store public key data via the PublicKeyCreator
// implementation.
func NewHandler(publicKeys PublicKeyCreator, logger hclog.Logger) *Handler {
	return &Handler{publicKeys: publicKeys, logger: logger.Named("signing")}
}

// Create a new public key for the client.
func (h *Handler) Create(ctx context.Context, payload *signingcmd.CreatePublicKey) error {
	if err := h.publicKeys.Create(ctx, payload.GetClientId(), payload.GetPublicKey()); err != nil {
		return fmt.Errorf("failed to create: %w", err)
	}

	h.logger.Debug("created public signing key", "client_id", payload.GetClientId())
	return nil
}
