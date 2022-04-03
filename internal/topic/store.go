package topic

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"

	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"

	"github.com/davidsbond/arrebato/internal/proto/arrebato/topic/v1"
)

type (
	// The BoltStore type is responsible for querying/mutating topic data within a boltdb database.
	BoltStore struct {
		db *bbolt.DB
	}
)

// NewBoltStore returns a new instance of the BoltStore type that will manage/query topic data in a boltdb database.
func NewBoltStore(db *bbolt.DB) *BoltStore {
	return &BoltStore{db: db}
}

var (
	// ErrTopicExists is the error given when attempting to create a new topic with the same name as an existing
	// topic.
	ErrTopicExists = errors.New("topic exists")

	// ErrNoTopic is the error given when attempting to delete/query a topic that does not exist.
	ErrNoTopic = errors.New("no topic")

	// ErrNoTopicInfo is the error given when querying a topic that has incomplete/invalid data within the state.
	ErrNoTopicInfo = errors.New("no topic info")
)

const (
	topicInfoKey = "info"
	topicsKey    = "topics"
	messagesKey  = "messages"
)

// Create a new topic. Returns ErrTopicExists if a topic with the same name already exists.
func (bs *BoltStore) Create(_ context.Context, t *topic.Topic) error {
	return bs.db.Update(func(tx *bbolt.Tx) error {
		topicsBucket, err := tx.CreateBucketIfNotExists([]byte(topicsKey))
		if err != nil {
			return fmt.Errorf("failed to open topics bucket: %w", err)
		}

		topicBucket, err := topicsBucket.CreateBucket([]byte(t.GetName()))
		switch {
		case errors.Is(err, bbolt.ErrBucketExists):
			return ErrTopicExists
		case err != nil:
			return fmt.Errorf("failed to create bucket: %w", err)
		}

		key := []byte(topicInfoKey)
		value, err := proto.Marshal(t)
		if err != nil {
			return fmt.Errorf("failed to marshal topic: %w", err)
		}

		if err = topicBucket.Put(key, value); err != nil {
			return fmt.Errorf("failed to store topic: %w", err)
		}

		// Prepare child buckets for messages and partitions.
		messagesBucket, err := topicBucket.CreateBucket([]byte(messagesKey))
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}

		for i := uint32(0); i < t.GetPartitions(); i++ {
			key = make([]byte, 8)
			binary.BigEndian.PutUint32(key, i)

			if _, err = messagesBucket.CreateBucket(key); err != nil {
				return fmt.Errorf("failed to create bucket: %w", err)
			}
		}

		return nil
	})
}

// Delete a named topic, returns ERrNoTopic if the topic does not exist.
func (bs *BoltStore) Delete(_ context.Context, t string) error {
	return bs.db.Update(func(tx *bbolt.Tx) error {
		topicsBucket, err := tx.CreateBucketIfNotExists([]byte(topicsKey))
		if err != nil {
			return fmt.Errorf("failed to open topics bucket: %w", err)
		}

		err = topicsBucket.DeleteBucket([]byte(t))
		switch {
		case errors.Is(err, bbolt.ErrBucketNotFound):
			return ErrNoTopic
		case err != nil:
			return fmt.Errorf("failed to delete bucket: %w", err)
		default:
			return nil
		}
	})
}

// Get information about a Topic by name. Returns ErrNoTopic if the topic does not exist.
func (bs *BoltStore) Get(ctx context.Context, name string) (*topic.Topic, error) {
	var t topic.Topic

	err := bs.db.View(func(tx *bbolt.Tx) error {
		topicsBucket := tx.Bucket([]byte(topicsKey))
		if topicsBucket == nil {
			return ErrNoTopic
		}

		topicBucket := topicsBucket.Bucket([]byte(name))
		if topicBucket == nil {
			return ErrNoTopic
		}

		val := topicBucket.Get([]byte(topicInfoKey))
		if val == nil {
			return ErrNoTopicInfo
		}

		if err := proto.Unmarshal(val, &t); err != nil {
			return fmt.Errorf("failed to unmarshal topic: %w", err)
		}

		return nil
	})

	return &t, err
}

// List all known topics. Returns ErrNoTopic if the store should contain information for a topic but doesn't.
func (bs *BoltStore) List(ctx context.Context) ([]*topic.Topic, error) {
	topics := make([]*topic.Topic, 0)
	err := bs.db.View(func(tx *bbolt.Tx) error {
		topicsBucket := tx.Bucket([]byte(topicsKey))
		if topicsBucket == nil {
			return nil
		}

		return topicsBucket.ForEach(func(k, v []byte) error {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			bucket := topicsBucket.Bucket(k)
			if bucket == nil {
				return nil
			}

			key := []byte(topicInfoKey)
			value := bucket.Get(key)
			if value == nil {
				return ErrNoTopicInfo
			}

			var t topic.Topic
			if err := proto.Unmarshal(value, &t); err != nil {
				return err
			}

			topics = append(topics, &t)
			return nil
		})
	})

	return topics, err
}
