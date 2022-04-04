package node_test

import (
	"context"

	"github.com/hashicorp/raft"

	"github.com/davidsbond/arrebato/internal/proto/arrebato/node/v1"
)

type (
	MockRaft struct {
		state  raft.RaftState
		config raft.Configuration
	}

	MockConfigurationFuture struct {
		raft.ConfigurationFuture

		config raft.Configuration
	}

	MockModifier struct {
		added *node.Node
		err   error
	}
)

func (mm *MockModifier) AllocateTopic(ctx context.Context, name string, topic string) error {
	return mm.err
}

func (mm *MockModifier) Add(ctx context.Context, node *node.Node) error {
	mm.added = node
	return mm.err
}

func (mm *MockModifier) Remove(ctx context.Context, name string) error {
	return mm.err
}

func (mm *MockRaft) GetConfiguration() raft.ConfigurationFuture {
	return &MockConfigurationFuture{config: mm.config}
}

func (mm *MockRaft) State() raft.RaftState {
	return mm.state
}

func (mm *MockConfigurationFuture) Error() error {
	return nil
}

func (mm *MockConfigurationFuture) Configuration() raft.Configuration {
	return mm.config
}
