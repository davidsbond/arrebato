package signing_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidsbond/arrebato/internal/signing"
	"github.com/davidsbond/arrebato/internal/testutil"
)

func TestBoltStore_Create(t *testing.T) {
	t.Parallel()

	db := testutil.BoltDB(t)
	ctx := testutil.Context(t)
	keys := signing.NewBoltStore(db)

	const clientID = "test"
	key := []byte("test")

	t.Run("It should save the public key", func(t *testing.T) {
		require.NoError(t, keys.Create(ctx, clientID, key))
	})

	t.Run("It should return the public key", func(t *testing.T) {
		actual, err := keys.Get(ctx, clientID)
		require.NoError(t, err)
		assert.EqualValues(t, key, actual)
	})
}

func TestBoltStore_Get(t *testing.T) {
	t.Parallel()

	db := testutil.BoltDB(t)
	ctx := testutil.Context(t)
	keys := signing.NewBoltStore(db)

	const clientID = "test"

	t.Run("It should return an error if a public key does not exist", func(t *testing.T) {
		_, err := keys.Get(ctx, clientID)
		require.Error(t, err)
	})
}
