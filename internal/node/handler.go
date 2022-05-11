package node

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	nodecmd "github.com/davidsbond/arrebato/internal/proto/arrebato/node/command/v1"
	"github.com/davidsbond/arrebato/internal/proto/arrebato/node/v1"
	"github.com/davidsbond/arrebato/internal/tracing"
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
		UnassignTopic(ctx context.Context, topicName string) error
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
	return tracing.WithinSpan(ctx, "Node.Add", func(ctx context.Context, span trace.Span) error {
		span.SetAttributes(
			attribute.String("node.name", payload.GetNode().GetName()),
		)

		if err := h.nodes.Create(ctx, payload.GetNode()); err != nil {
			return fmt.Errorf("failed to handle command: %w", err)
		}

		h.logger.Debug("node added", "name", payload.GetNode().GetName())
		return nil
	})
}

// Remove an existing node from the state.
func (h *Handler) Remove(ctx context.Context, payload *nodecmd.RemoveNode) error {
	return tracing.WithinSpan(ctx, "Node.Remove", func(ctx context.Context, span trace.Span) error {
		span.SetAttributes(
			attribute.String("node.name", payload.GetName()),
		)

		if err := h.nodes.Delete(ctx, payload.GetName()); err != nil {
			return fmt.Errorf("failed to handle command: %w", err)
		}

		h.logger.Debug("node removed", "name", payload.GetName())
		return nil
	})
}

// AssignTopic assigns a topic to a desired node.
func (h *Handler) AssignTopic(ctx context.Context, payload *nodecmd.AssignTopic) error {
	return tracing.WithinSpan(ctx, "Node.AssignTopic", func(ctx context.Context, span trace.Span) error {
		span.SetAttributes(
			attribute.String("node.name", payload.GetNodeName()),
			attribute.String("topic.name", payload.GetTopicName()),
		)

		if err := h.nodes.AssignTopic(ctx, payload.GetNodeName(), payload.GetTopicName()); err != nil {
			return fmt.Errorf("failed to handle command: %w", err)
		}

		h.logger.Debug("assigned topic", "node_name", payload.GetNodeName(), "topic_name", payload.GetTopicName())
		return nil
	})
}

// UnassignTopic unassigns a topic from any nodes it may be assigned to.
func (h *Handler) UnassignTopic(ctx context.Context, payload *nodecmd.UnassignTopic) error {
	return tracing.WithinSpan(ctx, "Node.UnassignTopic", func(ctx context.Context, span trace.Span) error {
		span.SetAttributes(
			attribute.String("node.name", payload.GetName()),
		)

		if err := h.nodes.UnassignTopic(ctx, payload.GetName()); err != nil {
			return fmt.Errorf("failed to handle command: %w", err)
		}

		h.logger.Debug("unassigned topic", "topic_name", payload.GetName())
		return nil
	})
}
