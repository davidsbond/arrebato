package node_test

import "github.com/hashicorp/raft"

type (
	MockRaft struct {
		state raft.RaftState
	}
)

func (mm *MockRaft) State() raft.RaftState {
	return mm.state
}
