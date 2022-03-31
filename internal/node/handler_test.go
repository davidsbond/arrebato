package node_test

import (
	"io"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	"github.com/davidsbond/arrebato/internal/node"
	nodecmd "github.com/davidsbond/arrebato/internal/proto/arrebato/node/command/v1"
	nodepb "github.com/davidsbond/arrebato/internal/proto/arrebato/node/v1"
	"github.com/davidsbond/arrebato/internal/testutil"
)

func TestHandler_Add(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name         string
		Command      *nodecmd.AddNode
		ExpectsError bool
		Error        error
	}{
		{
			Name: "It should add a new node",
			Command: &nodecmd.AddNode{
				Node: &nodepb.Node{
					Name: "node-0",
				},
			},
		},
		{
			Name:  "It should do nothing if the node already exists",
			Error: node.ErrNodeExists,
			Command: &nodecmd.AddNode{
				Node: &nodepb.Node{
					Name: "node-0",
				},
			},
		},
		{
			Name:         "It should propagate other errors",
			Error:        io.EOF,
			ExpectsError: true,
			Command: &nodecmd.AddNode{
				Node: &nodepb.Node{
					Name: "node-0",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := testutil.Context(t)
			store := &MockStore{err: tc.Error}

			err := node.NewHandler(store, hclog.NewNullLogger()).Add(ctx, tc.Command)
			if tc.ExpectsError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.True(t, proto.Equal(tc.Command.GetNode(), store.added))
		})
	}
}

func TestHandler_Remove(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name         string
		Command      *nodecmd.RemoveNode
		ExpectsError bool
		Error        error
	}{
		{
			Name: "It should remove an existing node",
			Command: &nodecmd.RemoveNode{
				Name: "node-0",
			},
		},
		{
			Name:  "It should do nothing if the node doesn't exist",
			Error: node.ErrNoNode,
			Command: &nodecmd.RemoveNode{
				Name: "node-0",
			},
		},
		{
			Name:         "It should propagate other errors",
			Error:        io.EOF,
			ExpectsError: true,
			Command: &nodecmd.RemoveNode{
				Name: "node-0",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := testutil.Context(t)
			store := &MockStore{err: tc.Error}

			err := node.NewHandler(store, hclog.NewNullLogger()).Remove(ctx, tc.Command)
			if tc.ExpectsError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
