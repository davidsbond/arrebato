package message

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	messagecmd "github.com/davidsbond/arrebato/internal/proto/arrebato/message/command/v1"
	"github.com/davidsbond/arrebato/internal/proto/arrebato/message/v1"
	"github.com/davidsbond/arrebato/internal/tracing"
)

type (
	// The Handler type is responsible for handling message commands and modifying the message state appropriately.
	Handler struct {
		logger   hclog.Logger
		messages Creator
	}

	// The Creator interface describes types that can create messages within a store.
	Creator interface {
		// Create should create a new message in the store and return the index of the message in the log.
		Create(ctx context.Context, m *message.Message) (uint64, error)
	}
)

// NewHandler returns a new instance of the Handler type that will store messages in the provided Creator implementation.
func NewHandler(messages Creator, logger hclog.Logger) *Handler {
	return &Handler{
		logger:   logger.Named("message"),
		messages: messages,
	}
}

// Create handles the messagecmd.CreateMessage command and creates a new message within the message store.
func (h *Handler) Create(ctx context.Context, cmd *messagecmd.CreateMessage) error {
	return tracing.WithinSpan(ctx, "Message.Create", func(ctx context.Context, span trace.Span) error {
		span.SetAttributes(
			attribute.String("topic.name", cmd.GetMessage().GetTopic()),
			attribute.String("sender.id", cmd.GetMessage().GetSender().GetId()),
		)

		index, err := h.messages.Create(ctx, cmd.GetMessage())
		if err != nil {
			return fmt.Errorf("failed to handle command: %w", err)
		}

		h.logger.Debug("message created",
			"index", index,
			"topic", cmd.GetMessage().GetTopic())

		return nil
	})
}
