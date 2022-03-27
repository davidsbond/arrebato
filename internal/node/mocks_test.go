package node_test

import "github.com/hashicorp/raft"

type (
	MockRaft struct {
		state  raft.RaftState
		config raft.Configuration
	}

	MockConfigurationFuture struct {
		raft.ConfigurationFuture

		config raft.Configuration
	}
)

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
