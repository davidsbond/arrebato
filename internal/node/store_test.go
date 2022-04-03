package node_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidsbond/arrebato/internal/node"
	nodepb "github.com/davidsbond/arrebato/internal/proto/arrebato/node/v1"
	"github.com/davidsbond/arrebato/internal/testutil"
)

func TestBoltStore_Add(t *testing.T) {
	t.Parallel()

	ctx := testutil.Context(t)
	db := testutil.BoltDB(t)

	store := node.NewBoltStore(db)

	t.Run("It should add a node", func(t *testing.T) {
		assert.NoError(t, store.Add(ctx, &nodepb.Node{
			Name: "node-0",
		}))
	})

	t.Run("It should return an error adding a duplicate node", func(t *testing.T) {
		err := store.Add(ctx, &nodepb.Node{
			Name: "node-0",
		})

		assert.True(t, errors.Is(err, node.ErrNodeExists))
	})
}

func TestBoltStore_Remove(t *testing.T) {
	t.Parallel()

	ctx := testutil.Context(t)
	db := testutil.BoltDB(t)

	store := node.NewBoltStore(db)

	t.Run("It should return an error if the node does not exist", func(t *testing.T) {
		err := store.Remove(ctx, "node-0")
		assert.True(t, errors.Is(err, node.ErrNoNode))
	})

	t.Run("It should remove an existing node", func(t *testing.T) {
		require.NoError(t, store.Add(ctx, &nodepb.Node{
			Name: "node-0",
		}))

		require.NoError(t, store.Remove(ctx, "node-0"))

		err := store.Remove(ctx, "node-0")
		assert.True(t, errors.Is(err, node.ErrNoNode))
	})
}

func TestBoltStore_LeastTopics(t *testing.T) {
	t.Parallel()

	ctx := testutil.Context(t)
	db := testutil.BoltDB(t)

	store := node.NewBoltStore(db)

	// Add nodes with some topics
	require.NoError(t, store.Add(ctx, &nodepb.Node{
		Name:   "node-0",
		Topics: []string{"topic-0", "topic-1"},
	}))

	require.NoError(t, store.Add(ctx, &nodepb.Node{
		Name:   "node-1",
		Topics: []string{"topic-2"},
	}))

	t.Run("It should return the node with the least topics", func(t *testing.T) {
		actual, err := store.LeastTopics(ctx)
		require.NoError(t, err)
		assert.EqualValues(t, "node-1", actual.GetName())
	})
}

func TestBoltStore_AllocateTopic(t *testing.T) {
	t.Parallel()

	ctx := testutil.Context(t)
	db := testutil.BoltDB(t)

	store := node.NewBoltStore(db)

	// Add a node to allocate a topic to
	require.NoError(t, store.Add(ctx, &nodepb.Node{
		Name: "node-0",
	}))

	t.Run("It should return an error if the node does not exist", func(t *testing.T) {
		err := store.AllocateTopic(ctx, "node-1", "topic-0")
		assert.True(t, errors.Is(err, node.ErrNoNode))
	})

	t.Run("It should allocate the topic to the existing node", func(t *testing.T) {
		require.NoError(t, store.AllocateTopic(ctx, "node-0", "topic-0"))
	})

	t.Run("It should do nothing if the topic is already allocated to that node", func(t *testing.T) {
		require.NoError(t, store.AllocateTopic(ctx, "node-0", "topic-0"))
	})
}
