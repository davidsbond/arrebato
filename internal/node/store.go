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
	// The BoltStore type is responsible for managing node data within a boltdb database.
	BoltStore struct {
		db *bbolt.DB
	}
)

const (
	nodesKey = "nodes"
)

var (
	// ErrNodeExists is the error given when an operation would overwrite the data for an existing node.
	ErrNodeExists = errors.New("node exists")
	// ErrNoNode is the error given when querying a node that does not exist.
	ErrNoNode = errors.New("no node")
)

// NewBoltStore returns a new instance of the BoltStore type that will store node data within the provided
// bbolt.DB instance.
func NewBoltStore(db *bbolt.DB) *BoltStore {
	return &BoltStore{db: db}
}

// Add a node to the store. Returns ErrNodeExists if a record already exists for the node.
func (bs *BoltStore) Add(ctx context.Context, node *node.Node) error {
	return bs.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(nodesKey))
		if err != nil {
			return fmt.Errorf("failed to open node bucket: %w", err)
		}

		key := []byte(node.GetName())
		if value := bucket.Get(key); len(value) > 0 {
			return ErrNodeExists
		}

		value, err := proto.Marshal(node)
		if err != nil {
			return fmt.Errorf("failed to marshal node: %w", err)
		}

		if err = bucket.Put(key, value); err != nil {
			return fmt.Errorf("failed to store node: %w", err)
		}

		return nil
	})
}

// Remove a node from the store. Returns ErrNoNode if a record does not exist for the node.
func (bs *BoltStore) Remove(ctx context.Context, name string) error {
	return bs.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(nodesKey))
		if err != nil {
			return fmt.Errorf("failed to open node bucket: %w", err)
		}

		key := []byte(name)
		if value := bucket.Get(key); len(value) == 0 {
			return ErrNoNode
		}

		if err = bucket.Delete(key); err != nil {
			return fmt.Errorf("failed to delete node: %w", err)
		}

		return nil
	})
}

// LeastTopics returns the node with the least number of allocated topics.
func (bs *BoltStore) LeastTopics(ctx context.Context) (*node.Node, error) {
	nodes := make([]*node.Node, 0)
	err := bs.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(nodesKey))
		if bucket == nil {
			return ErrNoNode
		}

		return bucket.ForEach(func(_, v []byte) error {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			var n node.Node
			if err := proto.Unmarshal(v, &n); err != nil {
				return fmt.Errorf("failed to unmarshal node: %w", err)
			}

			nodes = append(nodes, &n)
			return nil
		})
	})

	var selected *node.Node
	var lastCount int

	// Find the node with the least amount of allocated topics.
	for _, n := range nodes {
		if len(n.GetTopics()) < lastCount {
			selected = n
		}

		lastCount = len(n.GetTopics())
	}

	// If we couldn't find one, or all nodes had the same number of topics,
	// we'll just return the first one.
	if selected == nil {
		selected = nodes[0]
	}

	return selected, err
}

// AllocateTopic adds a topic to a node's allocated topics list. These are the topics that the node will permit
// consumers for.
func (bs *BoltStore) AllocateTopic(ctx context.Context, name, topic string) error {
	return bs.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(nodesKey))
		if bucket == nil {
			return ErrNoNode
		}

		key := []byte(name)
		value := bucket.Get(key)
		if len(value) == 0 {
			return ErrNoNode
		}

		var n node.Node
		if err := proto.Unmarshal(value, &n); err != nil {
			return fmt.Errorf("failed to unmarshal node: %w", err)
		}

		for _, t := range n.GetTopics() {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			if t != topic {
				continue
			}

			return nil
		}

		n.Topics = append(n.Topics, topic)
		value, err := proto.Marshal(&n)
		if err != nil {
			return fmt.Errorf("failed to marshal node: %w", err)
		}

		if err = bucket.Put(key, value); err != nil {
			return fmt.Errorf("failed to update node: %w", err)
		}

		return nil
	})
}
