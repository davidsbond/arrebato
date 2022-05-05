package node_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidsbond/arrebato/internal/node"
	nodepb "github.com/davidsbond/arrebato/internal/proto/arrebato/node/v1"
	"github.com/davidsbond/arrebato/internal/testutil"
)

func TestBoltStore_Create(t *testing.T) {
	t.Parallel()

	ctx := testutil.Context(t)
	db := testutil.BoltDB(t)
	store := node.NewBoltStore(db)

	t.Run("It should store a node in the database", func(t *testing.T) {
		record := &nodepb.Node{
			Name: "test",
		}

		require.NoError(t, store.Create(ctx, record))
	})
}

func TestBoltStore_Delete(t *testing.T) {
	t.Parallel()

	ctx := testutil.Context(t)
	db := testutil.BoltDB(t)
	store := node.NewBoltStore(db)

	// Create a record to remove
	require.NoError(t, store.Create(ctx, &nodepb.Node{
		Name: "test",
	}))

	t.Run("It should remove a node from the database", func(t *testing.T) {
		require.NoError(t, store.Delete(ctx, "test"))
	})

	t.Run("It should not return an error for a non-existent node", func(t *testing.T) {
		require.NoError(t, store.Delete(ctx, "test"))
	})
}

func TestBoltStore_AssignTopic(t *testing.T) {
	t.Parallel()

	ctx := testutil.Context(t)
	db := testutil.BoltDB(t)
	store := node.NewBoltStore(db)

	// Create some nodes to assign topics to
	require.NoError(t, store.Create(ctx, &nodepb.Node{
		Name: "node-0",
	}))
	require.NoError(t, store.Create(ctx, &nodepb.Node{
		Name: "node-1",
	}))

	const topicName = "test-topic"
	t.Run("It should assign the topic to a node", func(t *testing.T) {
		const nodeName = "node-0"

		require.NoError(t, store.AssignTopic(ctx, nodeName, topicName))

		nodes, err := store.List(ctx)
		require.NoError(t, err)
		for _, n := range nodes {
			if n.GetName() == nodeName {
				assert.Contains(t, n.GetTopics(), topicName)
				continue
			}

			assert.NotContains(t, n.GetTopics(), topicName)
		}
	})

	t.Run("It should reassign the topic to another node", func(t *testing.T) {
		const nodeName = "node-1"

		require.NoError(t, store.AssignTopic(ctx, nodeName, topicName))

		nodes, err := store.List(ctx)
		require.NoError(t, err)
		for _, n := range nodes {
			if n.GetName() == nodeName {
				assert.Contains(t, n.GetTopics(), topicName)
				continue
			}

			assert.NotContains(t, n.GetTopics(), topicName)
		}
	})
}

func TestBoltStore_GetTopicOwner(t *testing.T) {
	t.Parallel()

	ctx := testutil.Context(t)
	db := testutil.BoltDB(t)
	store := node.NewBoltStore(db)

	// Create some nodes to assign topics to
	require.NoError(t, store.Create(ctx, &nodepb.Node{
		Name: "node-0",
	}))
	require.NoError(t, store.Create(ctx, &nodepb.Node{
		Name: "node-1",
	}))

	// Assign a topic to a node
	const topicName = "topic-0"
	const nodeName = "node-0"
	require.NoError(t, store.AssignTopic(ctx, nodeName, topicName))

	t.Run("It should return the node assigned to the topic", func(t *testing.T) {
		result, err := store.GetTopicOwner(ctx, topicName)
		require.NoError(t, err)
		assert.EqualValues(t, nodeName, result.GetName())
	})

	t.Run("It should return an error if a node cannot be found for a topic", func(t *testing.T) {
		result, err := store.GetTopicOwner(ctx, "aaaaaaaaa")
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestBoltStore_UnassignTopic(t *testing.T) {
	t.Parallel()

	ctx := testutil.Context(t)
	db := testutil.BoltDB(t)
	store := node.NewBoltStore(db)

	// Create a node to assign topics to
	require.NoError(t, store.Create(ctx, &nodepb.Node{
		Name: "node-0",
	}))

	// Assign a topic to a node
	const topicName = "topic-0"
	const nodeName = "node-0"
	require.NoError(t, store.AssignTopic(ctx, nodeName, topicName))

	t.Run("It should unassign a topic from a node", func(t *testing.T) {
		require.NoError(t, store.UnassignTopic(ctx, topicName))
		_, err := store.GetTopicOwner(ctx, topicName)
		assert.Error(t, err)
	})
}
