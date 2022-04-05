package signing_test

import (
	"strconv"
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

func TestBoltStore_List(t *testing.T) {
	t.Parallel()

	db := testutil.BoltDB(t)
	ctx := testutil.Context(t)
	keys := signing.NewBoltStore(db)

	// Insert some keys to list
	for i := 0; i < 100; i++ {
		require.NoError(t, keys.Create(ctx, strconv.Itoa(i), []byte("test")))
	}

	t.Run("It should list all public keys", func(t *testing.T) {
		list, err := keys.List(ctx)
		require.NoError(t, err)
		assert.Len(t, list, 100)
		for _, k := range list {
			assert.NotEmpty(t, k.GetClientId())
			assert.NotEmpty(t, k.GetPublicKey())
		}
	})
}
