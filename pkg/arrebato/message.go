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
		Value  proto.Message
		Key    proto.Message
		Sender Sender
	}

	// The Sender type describes the client that produced a message.
	Sender struct {
		ID       string
		Verified bool
	}
)

func (m Message) toProto() (*message.Message, error) {
	value, err := anypb.New(m.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to convert value to any: %w", err)
	}

	var key *anypb.Any
	if m.Key != nil {
		key, err = anypb.New(m.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to convert key to any: %w", err)
		}
	}

	return &message.Message{
		Value: value,
		Key:   key,
	}, nil
}
