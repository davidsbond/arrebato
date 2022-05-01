package node_test

import (
	"testing"

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

	t.Run("It should not return an error for a non-existant node", func(t *testing.T) {
		require.NoError(t, store.Delete(ctx, "test"))
	})
}
