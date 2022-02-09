package server

import (
	"compress/gzip"
	"fmt"

	"github.com/hashicorp/raft"
	"go.etcd.io/bbolt"
)

type (
	// The Snapshot type is used to obtain a copy of the current data store managed by the raft log.
	Snapshot struct {
		store *bbolt.DB
	}
)

// Persist the contents of the data store to the raft.SnapshotSink implementation.
func (s *Snapshot) Persist(sink raft.SnapshotSink) error {
	return s.store.View(func(tx *bbolt.Tx) error {
		writer := gzip.NewWriter(sink)
		if _, err := tx.WriteTo(writer); err != nil {
			return fmt.Errorf("failed to write snapshot: %w", err)
		}

		return writer.Close()
	})
}

// Release does nothing.
func (s *Snapshot) Release() {}
