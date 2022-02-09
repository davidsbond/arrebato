// Package testutil provides utility functions that reduce duplicated code in tests.
package testutil

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.etcd.io/bbolt"
)

// BoltDB returns a new boltdb database in a temporary directory that is removed from-disk once the test is complete.
func BoltDB(t *testing.T) *bbolt.DB {
	t.Helper()

	dir, err := os.MkdirTemp(os.TempDir(), "arrebato")
	require.NoError(t, err)

	file := filepath.Join(dir, "arrebato.db")
	db, err := bbolt.Open(file, 0o755, bbolt.DefaultOptions)
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, db.Close())
		assert.NoError(t, os.RemoveAll(dir))
	})

	return db
}

// Context returns a new context.Context that is cancelled when the test is complete.
func Context(t *testing.T) context.Context {
	t.Helper()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	t.Cleanup(func() {
		cancel()
	})

	return ctx
}

// Any converts the provided proto.Message into an anypb.Any.
func Any(t *testing.T, m proto.Message) *anypb.Any {
	t.Helper()

	a, err := anypb.New(m)
	require.NoError(t, err)
	return a
}
