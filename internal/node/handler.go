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
	// the state within the Modifier.
	Handler struct {
		nodes  Modifier
		logger hclog.Logger
	}

	// The Modifier interface describes types that modifys node data.
	Modifier interface {
		// Add should add a new node to the store and return ErrNodeExists if a record with the same name already
		// exists.
		Add(ctx context.Context, node *node.Node) error

		// Remove should remove a node from the store and return ErrNoNode if a record with the node name does
		// not exist.
		Remove(ctx context.Context, name string) error

		// AllocateTopic should add a topic to the list of a node's allocated topics.
		AllocateTopic(ctx context.Context, name string, topic string) error
	}
)

// NewHandler returns a new instance of the Handler type that will persist node data within the provided
// Modifier implementation.
func NewHandler(modifier Modifier, logger hclog.Logger) *Handler {
	return &Handler{
		nodes:  modifier,
		logger: logger.Named("node"),
	}
}

// Add a new node to the store.
func (h *Handler) Add(ctx context.Context, cmd *nodecmd.AddNode) error {
	err := h.nodes.Add(ctx, cmd.GetNode())
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
	err := h.nodes.Remove(ctx, cmd.GetName())
	switch {
	case errors.Is(err, ErrNoNode):
		return nil
	case err != nil:
		return fmt.Errorf("failed to remove node: %w", err)
	}

	h.logger.Debug("removed node", "name", cmd.GetName())
	return nil
}

// AllocateTopic adds a topic to a node's allocated topic list.
func (h *Handler) AllocateTopic(ctx context.Context, cmd *nodecmd.AllocateTopic) error {
	if err := h.nodes.AllocateTopic(ctx, cmd.GetName(), cmd.GetTopic()); err != nil {
		return err
	}

	h.logger.Debug("allocated topic", "node", cmd.GetName(), "topic", cmd.GetTopic())
	return nil
}
