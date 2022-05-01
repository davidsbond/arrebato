package node

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"

	nodecmd "github.com/davidsbond/arrebato/internal/proto/arrebato/node/command/v1"
	"github.com/davidsbond/arrebato/internal/proto/arrebato/node/v1"
)

type (
	// The Handler type is responsible for handling raft commands regarding nodes and maintaining the node state.
	Handler struct {
		logger hclog.Logger
		nodes  Store
	}

	// The Store interface describes types that manage node state within the database.
	Store interface {
		Create(ctx context.Context, n *node.Node) error
		Delete(ctx context.Context, name string) error
		AssignTopic(ctx context.Context, nodeName, topicName string) error
	}
)

// NewHandler returns a new instance of the Handler type that will maintain node state via the Store implementation.
func NewHandler(nodes Store, logger hclog.Logger) *Handler {
	return &Handler{
		logger: logger.Named("node"),
		nodes:  nodes,
	}
}

// Add a new node to the state.
func (h *Handler) Add(ctx context.Context, payload *nodecmd.AddNode) error {
	if err := h.nodes.Create(ctx, payload.GetNode()); err != nil {
		return fmt.Errorf("failed to handle command: %w", err)
	}

	h.logger.Debug("node added", "name", payload.GetNode().GetName())
	return nil
}

// Remove an existing node from the state.
func (h *Handler) Remove(ctx context.Context, payload *nodecmd.RemoveNode) error {
	if err := h.nodes.Delete(ctx, payload.GetName()); err != nil {
		return fmt.Errorf("failed to handle command: %w", err)
	}

	h.logger.Debug("node removed", "name", payload.GetName())
	return nil
}

// AssignTopic assigns a topic to a desired node.
func (h *Handler) AssignTopic(ctx context.Context, payload *nodecmd.AssignTopic) error {
	if err := h.nodes.AssignTopic(ctx, payload.GetNodeName(), payload.GetTopicName()); err != nil {
		return fmt.Errorf("failed to handle command: %w", err)
	}

	h.logger.Debug("assigned topic", "node_name", payload.GetNodeName(), "topic_name", payload.GetTopicName())
	return nil
}