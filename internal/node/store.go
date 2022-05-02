package node

import (
	"context"
	"errors"
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

// ErrNodeExists is the error given when attempting to create a node that already exists in the store.
var ErrNodeExists = errors.New("node exists")

// Create a record for the node. Returns ErrNodeExists if a record already exists for the given node.
func (bs *BoltStore) Create(ctx context.Context, n *node.Node) error {
	return bs.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(nodeKey))
		if err != nil {
			return fmt.Errorf("failed to open node bucket: %w", err)
		}

		key := []byte(n.GetName())
		if value := bucket.Get(key); value != nil {
			return ErrNodeExists
		}

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

// List all nodes in the store.
func (bs *BoltStore) List(ctx context.Context) ([]*node.Node, error) {
	results := make([]*node.Node, 0)
	err := bs.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(nodeKey))
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			var n node.Node
			if err := proto.Unmarshal(v, &n); err != nil {
				return fmt.Errorf("failed to unmarshal node: %w", err)
			}

			results = append(results, &n)
			return nil
		})
	})

	return results, err
}

// AssignTopic adds the provided topic name to a node's list of topics. It also checks other node records and removes
// the topic from their list if present.
func (bs *BoltStore) AssignTopic(ctx context.Context, nodeName string, topicName string) error {
	return bs.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(nodeKey))
		if err != nil {
			return fmt.Errorf("failed to open node bucket: %w", err)
		}

		return bucket.ForEach(func(k, v []byte) error {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			var n node.Node
			if err = proto.Unmarshal(v, &n); err != nil {
				return fmt.Errorf("failed to unmarshal node: %w", err)
			}

			// If this is the desired node, ensure the topic is not already assigned to it, append the
			// new topic name and save the updated record.
			if n.GetName() == nodeName {
				for _, tp := range n.GetTopics() {
					if tp == topicName {
						return nil
					}
				}

				n.Topics = append(n.Topics, topicName)
				value, err := proto.Marshal(&n)
				if err != nil {
					return fmt.Errorf("failed to marshal node: %w", err)
				}

				return bucket.Put([]byte(n.GetName()), value)
			}

			// Remove any references to the topic from other node records
			for i, tp := range n.GetTopics() {
				if tp != topicName {
					continue
				}

				n.Topics = append(n.Topics[:i], n.Topics[i+1:]...)
				value, err := proto.Marshal(&n)
				if err != nil {
					return fmt.Errorf("failed to marshal node: %w", err)
				}

				return bucket.Put([]byte(n.GetName()), value)
			}

			return nil
		})
	})
}
