// Package acl provides all functionality within arrebato regarding access-control lists. This includes both gRPC, raft
// and data store interactions.
package acl

import (
	"context"
	"fmt"

	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"

	"github.com/davidsbond/arrebato/internal/proto/arrebato/acl/v1"
)

type (
	// The BoltStore type is responsible for querying/mutating ACL data within a boltdb database.
	BoltStore struct {
		db *bbolt.DB
	}
)

// NewBoltStore returns a new instance of the BoltStore type that will manage/query ACL data in a boltdb database.
func NewBoltStore(db *bbolt.DB) *BoltStore {
	return &BoltStore{db: db}
}

const (
	topicsKey = "topics"
	aclKey    = "acl"
)

// Set the ACL to the one provided.
func (bs *BoltStore) Set(_ context.Context, a *acl.ACL) error {
	value, err := proto.Marshal(a)
	if err != nil {
		return fmt.Errorf("failed to marshal acl: %w", err)
	}

	return bs.db.Update(func(tx *bbolt.Tx) error {
		topics, err := tx.CreateBucketIfNotExists([]byte(topicsKey))
		if err != nil {
			return fmt.Errorf("failed to open topics bucket: %w", err)
		}

		return topics.Put([]byte(aclKey), value)
	})
}

// Allowed returns a boolean value indicating if the client has the given permission on a topic. This method returns true
// in scenarios where an ACL has yet to be created.
func (bs *BoltStore) Allowed(ctx context.Context, topic, client string, permission acl.Permission) (bool, error) {
	var allowed bool
	err := bs.db.View(func(tx *bbolt.Tx) error {
		topics := tx.Bucket([]byte(topicsKey))

		// No topic bucket means no ACLs, so we allow everything until the first ACL rules are in.
		if topics == nil {
			allowed = true
			return nil
		}

		// If no ACL has been created yet, we assume we can do anything we want.
		value := topics.Get([]byte(aclKey))
		if value == nil {
			allowed = true
			return nil
		}

		var a acl.ACL
		if err := proto.Unmarshal(value, &a); err != nil {
			return fmt.Errorf("failed to unmarshal acl: %w", err)
		}

		for _, entry := range a.GetEntries() {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			if entry.GetTopic() != topic || entry.GetClient() != client {
				continue
			}

			for _, perm := range entry.GetPermissions() {
				if perm != permission {
					continue
				}

				allowed = true
				return nil
			}
		}

		return nil
	})

	return allowed, err
}
