package node

import (
	"context"
	"fmt"

	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"

	"github.com/davidsbond/arrebato/internal/proto/arrebato/node/v1"
)

type (
	// The BoltStore type is responsible for managing node state within a boltdb database.
	BoltStore struct {
		db *bbolt.DB
	}
)

const nodeKey = "node"

// NewBoltStore returns a new instance of the BoltStore type that will store node state within the bbolt.DB instance.
func NewBoltStore(db *bbolt.DB) *BoltStore {
	return &BoltStore{db: db}
}

// Create a record for the node.
func (bs *BoltStore) Create(ctx context.Context, n *node.Node) error {
	return bs.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(nodeKey))
		if err != nil {
			return fmt.Errorf("failed to open node bucket: %w", err)
		}

		key := []byte(n.GetName())
		value, err := proto.Marshal(n)
		if err != nil {
			return fmt.Errorf("failed to marshal node: %w", err)
		}

		if err = bucket.Put(key, value); err != nil {
			return fmt.Errorf("failed to store node: %w", err)
		}

		return nil
	})
}

// Delete the record of a named node. If the node does not exist, nothing happens and a nil error is returned.
func (bs *BoltStore) Delete(ctx context.Context, name string) error {
	return bs.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(nodeKey))
		if err != nil {
			return fmt.Errorf("failed to open node bucket: %w", err)
		}

		if err = bucket.Delete([]byte(name)); err != nil {
			return fmt.Errorf("failed to delete node: %w", err)
		}

		return nil
	})
}
