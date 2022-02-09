package arrebato

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/davidsbond/arrebato/internal/proto/arrebato/message/v1"
)

type (
	// The Message type describes a message that is produced/consumed by arrebato clients.
	Message struct {
		Payload proto.Message
	}
)

func (m Message) toProto() (*message.Message, error) {
	anyPayload, err := anypb.New(m.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to convert payload to any: %w", err)
	}

	return &message.Message{
		Payload: anyPayload,
	}, nil
}
