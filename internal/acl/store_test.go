package acl_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidsbond/arrebato/internal/acl"
	aclpb "github.com/davidsbond/arrebato/internal/proto/arrebato/acl/v1"
	"github.com/davidsbond/arrebato/internal/testutil"
)

func TestBoltStore_Allowed(t *testing.T) {
	t.Parallel()

	ctx := testutil.Context(t)
	db := testutil.BoltDB(t)

	// Create an ACL to perform tests against
	acls := acl.NewBoltStore(db)
	require.NoError(t, acls.Set(ctx, &aclpb.ACL{
		Entries: []*aclpb.Entry{
			{
				Topic:  "test-topic",
				Client: "test-client",
				Permissions: []aclpb.Permission{
					aclpb.Permission_PERMISSION_CONSUME,
				},
			},
		},
	}))

	t.Run("It should return false if the client does not have permission on the topic", func(t *testing.T) {
		allowed, err := acls.Allowed(ctx, "test-topic", "test-client", aclpb.Permission_PERMISSION_PRODUCE)
		assert.NoError(t, err)
		assert.False(t, allowed)
	})

	t.Run("It should return true if the client does have permission on the topic", func(t *testing.T) {
		allowed, err := acls.Allowed(ctx, "test-topic", "test-client", aclpb.Permission_PERMISSION_CONSUME)
		assert.NoError(t, err)
		assert.True(t, allowed)
	})
}
