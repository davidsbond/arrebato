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
		Error        error
		ExpectsError bool
		Command      *nodecmd.AddNode
		Expected     *nodepb.Node
	}{
		{
			Name: "It should add the node from the command",
			Command: &nodecmd.AddNode{
				Node: &nodepb.Node{
					Name: "test",
				},
			},
			Expected: &nodepb.Node{
				Name: "test",
			},
		},
		{
			Name:         "It should propagate errors",
			Error:        io.EOF,
			ExpectsError: true,
			Command: &nodecmd.AddNode{
				Node: &nodepb.Node{
					Name: "test",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			store := &MockStore{err: tc.Error}
			ctx := testutil.Context(t)

			err := node.NewHandler(store, hclog.NewNullLogger()).Add(ctx, tc.Command)
			if tc.ExpectsError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.True(t, proto.Equal(tc.Expected, store.saved))
		})
	}
}

func TestHandler_Remove(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name         string
		Error        error
		ExpectsError bool
		Command      *nodecmd.RemoveNode
	}{
		{
			Name: "It should add the node from the command",
			Command: &nodecmd.RemoveNode{
				Name: "test",
			},
		},
		{
			Name:         "It should propagate errors",
			Error:        io.EOF,
			ExpectsError: true,
			Command: &nodecmd.RemoveNode{
				Name: "test",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			store := &MockStore{err: tc.Error}
			ctx := testutil.Context(t)

			err := node.NewHandler(store, hclog.NewNullLogger()).Remove(ctx, tc.Command)
			if tc.ExpectsError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
