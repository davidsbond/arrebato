package message

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"

	"github.com/davidsbond/arrebato/internal/proto/arrebato/message/v1"
	"github.com/davidsbond/arrebato/internal/proto/arrebato/topic/v1"
)

type (
	// The BoltStore type is responsible for persisting/querying message data from a boltdb instance.
	BoltStore struct {
		db *bbolt.DB
	}

	// ReadFunc is a function that is invoked whenever a message is read from a topic.
	ReadFunc func(ctx context.Context, m *message.Message) error
)

// NewBoltStore returns a new instance of the BoltStore type that uses the provided bbolt.DB instance for
// persistence.
func NewBoltStore(db *bbolt.DB) *BoltStore {
	return &BoltStore{db: db}
}

var (
	// ErrNoTopic is the error given when attempting to read/write messages on a topic that does not exist.
	ErrNoTopic = errors.New("no topic")

	// ErrNoMessages is the error given when attempting to read messages from a topic that contains none.
	ErrNoMessages = errors.New("no messages")

	// ErrNoTopicInfo is the error given when a topic is found without any data regarding its retention
	// period. A topic should never be in this state and probably needs recreating.
	ErrNoTopicInfo = errors.New("no topic info")
)

const (
	topicInfoKey = "info"
	messagesKey  = "messages"
	topicsKey    = "topics"
)

// Create a new Message within the store, returning its index.
func (bs *BoltStore) Create(_ context.Context, m *message.Message) (uint64, error) {
	var index uint64

	err := bs.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(topicsKey))
		if bucket == nil {
			return ErrNoTopic
		}

		topicBucket := bucket.Bucket([]byte(m.GetTopic()))
		if topicBucket == nil {
			return ErrNoTopic
		}

		messages, err := topicBucket.CreateBucketIfNotExists([]byte(messagesKey))
		if err != nil {
			return fmt.Errorf("failed to open message bucket: %w", err)
		}

		partitionKey := make([]byte, 8)
		binary.BigEndian.PutUint32(partitionKey, m.GetPartition())

		partition, err := messages.CreateBucketIfNotExists(partitionKey)
		if err != nil {
			return fmt.Errorf("failed to open partition bucket: %w", err)
		}

		// We leverage BoltDB's ability to generate sequences so that we can have an ordinal number representing the
		// position in the log of the message. This will allow subscribers to read all messages from particular points
		// in time.
		index, err = partition.NextSequence()
		if err != nil {
			return fmt.Errorf("failed to generate index: %w", err)
		}

		// Boltdb sequence starts at 1, but I prefer zero, so we always decrement by one.
		index--
		m.Index = index

		// We store an encoded representation of the index so that keys are lexicographically sorted.
		key := make([]byte, 8)
		binary.BigEndian.PutUint64(key, index)

		value, err := proto.Marshal(m)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}

		if err = partition.Put(key, value); err != nil {
			return fmt.Errorf("failed to store message: %w", err)
		}

		return nil
	})

	return index, err
}

// Read messages from the desired topic, starting at the desired index. Each message triggers an invocation of the
// provided ReadFunc. This method blocks until the end of the messages is reached, the provided context is cancelled
// or the ReadFunc returns an error.
func (bs *BoltStore) Read(ctx context.Context, topic string, partition uint32, startIndex uint64, fn ReadFunc) error {
	return bs.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(topicsKey))
		if bucket == nil {
			return ErrNoTopic
		}

		topicBucket := bucket.Bucket([]byte(topic))
		if topicBucket == nil {
			return ErrNoTopic
		}

		messagesBucket := topicBucket.Bucket([]byte(messagesKey))
		if messagesBucket == nil {
			return ErrNoMessages
		}

		partitionKey := make([]byte, 8)
		binary.BigEndian.PutUint32(partitionKey, partition)

		partitionBucket := messagesBucket.Bucket(partitionKey)
		if partitionBucket == nil {
			return ErrNoMessages
		}

		cursor := partitionBucket.Cursor()

		// Keys are stored lexicographically as encoded uint64s, so we can iterate over them.
		start := make([]byte, 8)
		binary.BigEndian.PutUint64(start, startIndex)

		for k, v := cursor.Seek(start); k != nil; k, v = cursor.Next() {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			var m message.Message
			if err := proto.Unmarshal(v, &m); err != nil {
				return fmt.Errorf("failed to unmarshal message: %w", err)
			}

			if err := fn(ctx, &m); err != nil {
				return err
			}
		}

		return nil
	})
}

// Prune all messages on a topic created before the given timestamp.
func (bs *BoltStore) Prune(ctx context.Context, topicName string, before time.Time) (uint64, error) {
	var count uint64

	err := bs.db.Update(func(tx *bbolt.Tx) error {
		topicsBucket := tx.Bucket([]byte(topicsKey))
		if topicsBucket == nil {
			return ErrNoTopic
		}

		topicBucket := topicsBucket.Bucket([]byte(topicName))
		if topicBucket == nil {
			return ErrNoTopic
		}

		data := topicBucket.Get([]byte(topicInfoKey))
		if data == nil {
			return ErrNoTopicInfo
		}

		var t topic.Topic
		if err := proto.Unmarshal(data, &t); err != nil {
			return fmt.Errorf("failed to unmarshal topic: %w", err)
		}

		messages := topicBucket.Bucket([]byte(messagesKey))
		if messages == nil {
			// This bucket may not exist if no messages have yet to be produced on the topic.
			return nil
		}

		// Iterate over each partition and remove messages.
		return messages.ForEach(func(k, v []byte) error {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			partitionBucket := messages.Bucket(k)
			if partitionBucket == nil {
				return nil
			}

			cursor := partitionBucket.Cursor()

			// Keys are stored lexicographically as encoded uint64s, so we can iterate over them using a cursor,
			// so we don't have to use ForEach.
			start := make([]byte, 8)
			binary.BigEndian.PutUint64(start, 0)

			for k, v := cursor.Seek(start); k != nil; k, v = cursor.Next() {
				if ctx.Err() != nil {
					return ctx.Err()
				}

				var m message.Message
				if err := proto.Unmarshal(v, &m); err != nil {
					return fmt.Errorf("failed to unmarshal message: %w", err)
				}

				// Once we've reached the first message that does not fall outside the retention period, we can stop
				// iterating as we can safely assume all other messages are also within the retention period.
				if m.GetTimestamp().AsTime().After(before) {
					break
				}

				if err := cursor.Delete(); err != nil {
					return err
				}

				count++
			}

			return nil
		})
	})

	return count, err
}

// Indexes returns the current message index of every known topic.
func (bs *BoltStore) Indexes(ctx context.Context) (map[string]map[uint32]uint64, error) {
	result := make(map[string]map[uint32]uint64)

	err := bs.db.View(func(tx *bbolt.Tx) error {
		topicsBucket := tx.Bucket([]byte(topicsKey))
		if topicsBucket == nil {
			return nil
		}

		return topicsBucket.ForEach(func(topicKey, v []byte) error {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			topicBucket := topicsBucket.Bucket(topicKey)
			if topicBucket == nil {
				return nil
			}

			messages := topicBucket.Bucket([]byte(messagesKey))
			partitions := make(map[uint32]uint64)
			result[string(topicKey)] = partitions

			return messages.ForEach(func(partitionKey, v []byte) error {
				partition := messages.Bucket(partitionKey)
				if partition == nil {
					return nil
				}

				partitions[binary.BigEndian.Uint32(partitionKey)] = partition.Sequence()
				return nil
			})
		})
	})

	return result, err
}

// Counts returns the number of messages available within each topic.
func (bs *BoltStore) Counts(ctx context.Context) (map[string]uint64, error) {
	result := make(map[string]uint64)

	err := bs.db.View(func(tx *bbolt.Tx) error {
		topicsBucket := tx.Bucket([]byte(topicsKey))
		if topicsBucket == nil {
			return nil
		}

		return topicsBucket.ForEach(func(k, v []byte) error {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			topicBucket := topicsBucket.Bucket(k)
			if topicBucket == nil {
				return nil
			}

			messages := topicBucket.Bucket([]byte(messagesKey))
			if messages == nil {
				result[string(k)] = 0
				return nil
			}

			var total uint64
			err := messages.ForEach(func(partitionKey, v []byte) error {
				partition := messages.Bucket(partitionKey)
				if partition == nil {
					return nil
				}

				total += uint64(partition.Stats().KeyN)
				return nil
			})
			if err != nil {
				return err
			}

			result[string(k)] = total
			return nil
		})
	})

	return result, err
}
