package node

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-hclog"

	nodecmd "github.com/davidsbond/arrebato/internal/proto/arrebato/node/command/v1"
	"github.com/davidsbond/arrebato/internal/proto/arrebato/node/v1"
)

type (
	// The Handler type is responsible for handling inbound commands relating to nodes from the raft log and updating
	// the state within the Store.
	Handler struct {
		store  Store
		logger hclog.Logger
	}

	// The Store interface describes types that persist node data.
	Store interface {
		// Add should add a new node to the store and return ErrNodeExists if a record with the same name already
		// exists.
		Add(ctx context.Context, node *node.Node) error

		// Remove should remove a node from the store and return ErrNoNode if a record with the node name does
		// not exist.
		Remove(ctx context.Context, name string) error
	}
)

// NewHandler returns a new instance of the Handler type that will persist node data within the provided
// Store implementation.
func NewHandler(store Store, logger hclog.Logger) *Handler {
	return &Handler{
		store:  store,
		logger: logger.Named("node"),
	}
}

// Add a new node to the store.
func (h *Handler) Add(ctx context.Context, cmd *nodecmd.AddNode) error {
	err := h.store.Add(ctx, cmd.GetNode())
	switch {
	case errors.Is(err, ErrNodeExists):
		return nil
	case err != nil:
		return fmt.Errorf("failed to store node: %w", err)
	}

	h.logger.Debug("added node", "name", cmd.GetNode().GetName())
	return nil
}

// Remove a node from the store.
func (h *Handler) Remove(ctx context.Context, cmd *nodecmd.RemoveNode) error {
	err := h.store.Remove(ctx, cmd.GetName())
	switch {
	case errors.Is(err, ErrNoNode):
		return nil
	case err != nil:
		return fmt.Errorf("failed to remove node: %w", err)
	}

	h.logger.Debug("removed node", "name", cmd.GetName())
	return nil
}
