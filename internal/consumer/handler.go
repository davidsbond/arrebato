package consumer

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"

	consumercmd "github.com/davidsbond/arrebato/internal/proto/arrebato/consumer/command/v1"
	"github.com/davidsbond/arrebato/internal/proto/arrebato/consumer/v1"
)

type (
	// The Handler type is responsible for handling commands sent to the server regarding consumer state.
	Handler struct {
		manager Manager
		logger  hclog.Logger
	}

	// The Manager interface describes types that can manage consumer state.
	Manager interface {
		SetTopicIndex(ctx context.Context, c *consumer.TopicIndex) error
	}
)

// NewHandler returns a new instance of the Handler type that will handle inbound commands regarding consumers.
func NewHandler(manager Manager, logger hclog.Logger) *Handler {
	return &Handler{
		manager: manager,
		logger:  logger.Named("consumer"),
	}
}

// SetTopicIndex handles a command that modifies the current index for a consumer on a topic.
func (h *Handler) SetTopicIndex(ctx context.Context, cmd *consumercmd.SetTopicIndex) error {
	if err := h.manager.SetTopicIndex(ctx, cmd.GetTopicIndex()); err != nil {
		return fmt.Errorf("failed to set topic index: %w", err)
	}

	h.logger.Debug("set topic index",
		"topic", cmd.GetTopicIndex().GetTopic(),
		"index", cmd.GetTopicIndex().GetIndex(),
		"consumer_id", cmd.GetTopicIndex().GetConsumerId())

	return nil
}
