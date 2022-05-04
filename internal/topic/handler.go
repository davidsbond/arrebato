package topic

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"

	topiccmd "github.com/davidsbond/arrebato/internal/proto/arrebato/topic/command/v1"
	"github.com/davidsbond/arrebato/internal/proto/arrebato/topic/v1"
)

type (
	// The Handler type is responsible for handling Topic related commands passed to an Executor implementation and
	// modifying the state of topic data appropriately.
	Handler struct {
		topics Manager
		logger hclog.Logger
	}

	// The Manager interface describes types responsible for managing the state of Topics within a data store.
	Manager interface {
		Create(ctx context.Context, t *topic.Topic) error
		Delete(ctx context.Context, t string) error
	}
)

// NewHandler returns a new instance of the Handler type that will modify topic data via the provided Manager
// implementation.
func NewHandler(manager Manager, logger hclog.Logger) *Handler {
	return &Handler{topics: manager, logger: logger.Named("topic")}
}

// Create the topic described in the command.
func (h *Handler) Create(ctx context.Context, cmd *topiccmd.CreateTopic) error {
	if err := h.topics.Create(ctx, cmd.GetTopic()); err != nil {
		return fmt.Errorf("failed to create topic %s: %w", cmd.GetTopic().GetName(), err)
	}

	h.logger.Debug("topic created",
		"name", cmd.GetTopic().GetName(),
		"message_retention_period", cmd.GetTopic().GetMessageRetentionPeriod().AsDuration(),
		"consumer_retention_period", cmd.GetTopic().GetConsumerRetentionPeriod().AsDuration(),
		"require_verified_messages", cmd.GetTopic().GetRequireVerifiedMessages(),
	)

	return nil
}

// Delete the topic described in the command.
func (h *Handler) Delete(ctx context.Context, cmd *topiccmd.DeleteTopic) error {
	if err := h.topics.Delete(ctx, cmd.GetName()); err != nil {
		return fmt.Errorf("failed to delete topic %s: %w", cmd.GetName(), err)
	}

	h.logger.Debug("topic deleted", "name", cmd.GetName())
	return nil
}
