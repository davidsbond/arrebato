package acl

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"go.opentelemetry.io/otel/trace"

	aclcmd "github.com/davidsbond/arrebato/internal/proto/arrebato/acl/command/v1"
	"github.com/davidsbond/arrebato/internal/proto/arrebato/acl/v1"
	"github.com/davidsbond/arrebato/internal/tracing"
)

type (
	// The Handler type is responsible for handling commands sent to the server regarding ACL state.
	Handler struct {
		acls   Setter
		logger hclog.Logger
	}

	// The Setter interface describes types that can set ACL state.
	Setter interface {
		Set(ctx context.Context, a *acl.ACL) error
	}
)

// NewHandler returns a new instance of the Handler type that will handle inbound commands regarding ACLs.
func NewHandler(setter Setter, logger hclog.Logger) *Handler {
	return &Handler{
		acls:   setter,
		logger: logger.Named("acl"),
	}
}

// Set handles a command that modifies the current ACL state.
func (h *Handler) Set(ctx context.Context, cmd *aclcmd.SetACL) error {
	return tracing.WithinSpan(ctx, "ACL.Set", func(ctx context.Context, span trace.Span) error {
		if err := h.acls.Set(ctx, cmd.GetAcl()); err != nil {
			return fmt.Errorf("failed to set acl: %w", err)
		}

		h.logger.Debug("updated acl", "entries", len(cmd.GetAcl().GetEntries()))

		return nil
	})
}
