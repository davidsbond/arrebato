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

	MockStore struct {
		saved *node.Node
		err   error
	}
)

func (mm *MockStore) AssignTopic(ctx context.Context, nodeName, topicName string) error {
	return mm.err
}

func (mm *MockStore) Create(ctx context.Context, n *node.Node) error {
	mm.saved = n
	return mm.err
}

func (mm *MockStore) Delete(ctx context.Context, name string) error {
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
