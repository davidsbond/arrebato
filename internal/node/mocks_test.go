package node_test

import (
	"context"

	"github.com/hashicorp/raft"

	"github.com/davidsbond/arrebato/internal/proto/arrebato/node/v1"
)

type (
	MockRaft struct {
		state raft.RaftState
	}

	MockStore struct {
		saved *node.Node
		err   error
	}

	MockLister struct {
		nodes []*node.Node
		err   error
	}
)

func (mm *MockStore) UnassignTopic(ctx context.Context, topicName string) error {
	return mm.err
}

func (mm *MockLister) List(ctx context.Context) ([]*node.Node, error) {
	return mm.nodes, mm.err
}

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

func (mm *MockRaft) State() raft.RaftState {
	return mm.state
}
