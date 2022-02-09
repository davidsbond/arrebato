package command_test

import (
	"time"

	"github.com/hashicorp/raft"
)

type (
	MockApplier struct {
		err     error
		applied []byte
		state   raft.RaftState
	}

	MockApplyFuture struct {
		err error
	}
)

func (mm *MockApplyFuture) Error() error {
	return mm.err
}

func (mm *MockApplyFuture) Index() uint64 {
	return 0
}

func (mm *MockApplyFuture) Response() interface{} {
	return mm.err
}

func (mm *MockApplier) ApplyLog(log raft.Log, _ time.Duration) raft.ApplyFuture {
	mm.applied = log.Data
	return &MockApplyFuture{err: mm.err}
}

func (mm *MockApplier) State() raft.RaftState {
	return mm.state
}
