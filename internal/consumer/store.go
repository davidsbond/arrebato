package consumer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"

	"github.com/davidsbond/arrebato/internal/proto/arrebato/consumer/v1"
)

type (
	// The BoltStore type is responsible for maintaining the state of consumers within a boltdb database.
	BoltStore struct {
		db *bbolt.DB
	}
)

// NewBoltStore returns a new instance of the BoltStore type that will persist and query consumer state from a boltdb
// instance.
func NewBoltStore(db *bbolt.DB) *BoltStore {
	return &BoltStore{db: db}
}

// ErrNoTopic is the error given when attempting to perform an action against a topic that does not exist.
var ErrNoTopic = errors.New("no topic")

const (
	consumersKey = "consumers"
	topicsKey    = "topics"
)

// SetTopicIndex sets the current index on a topic for a consumer. Returns ErrNoTopic if the topic does not exist.
func (bs *BoltStore) SetTopicIndex(_ context.Context, c *consumer.TopicIndex) error {
	return bs.db.Update(func(tx *bbolt.Tx) error {
		topicsBucket := tx.Bucket([]byte(topicsKey))
		if topicsBucket == nil {
			return ErrNoTopic
		}

		topicBucket := topicsBucket.Bucket([]byte(c.GetTopic()))
		if topicBucket == nil {
			return ErrNoTopic
		}

		consumers, err := topicBucket.CreateBucketIfNotExists([]byte(consumersKey))
		if err != nil {
			return fmt.Errorf("failed to open consumer bucket: %w", err)
		}

		key := []byte(c.GetConsumerId())
		value, err := proto.Marshal(c)
		if err != nil {
			return fmt.Errorf("failed to marshal topic index: %w", err)
		}

		if err = consumers.Put(key, value); err != nil {
			return fmt.Errorf("failed to store index: %w", err)
		}

		return nil
	})
}

// GetTopicIndex returns the current index in a topic for a consumer identifier. Returns a zero index if the consumer
// does not exist. Or ErrNoTopic if the topic does not exist.
func (bs *BoltStore) GetTopicIndex(_ context.Context, topic, consumerID string) (*consumer.TopicIndex, error) {
	index := consumer.TopicIndex{
		Topic:      topic,
		ConsumerId: consumerID,
	}

	err := bs.db.View(func(tx *bbolt.Tx) error {
		topicsBucket := tx.Bucket([]byte(topicsKey))
		if topicsBucket == nil {
			return ErrNoTopic
		}

		topicBucket := topicsBucket.Bucket([]byte(topic))
		if topicBucket == nil {
			return ErrNoTopic
		}

		consumers := topicBucket.Bucket([]byte(consumersKey))
		if consumers == nil {
			return nil
		}

		key := []byte(consumerID)
		value := consumers.Get(key)
		if value == nil {
			return nil
		}

		if err := proto.Unmarshal(value, &index); err != nil {
			return fmt.Errorf("failed to unmarshal topic index data: %w", err)
		}

		return nil
	})

	return &index, err
}

// Indexes returns the current consumer index of every consumer for every topic. The initial map is keyed by topic
// name, the nested map is keyed by consumer identifier.
func (bs *BoltStore) Indexes(ctx context.Context) (map[string]map[string]uint64, error) {
	result := make(map[string]map[string]uint64)

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

			indexes := make(map[string]uint64)
			result[string(k)] = indexes

			consumers := topicBucket.Bucket([]byte(consumersKey))
			if consumers == nil {
				return nil
			}

			return consumers.ForEach(func(k, v []byte) error {
				if ctx.Err() != nil {
					return ctx.Err()
				}

				var index consumer.TopicIndex
				if err := proto.Unmarshal(v, &index); err != nil {
					return fmt.Errorf("failed to unmarshal topic index data: %w", err)
				}

				indexes[index.GetConsumerId()] = index.GetIndex()
				return nil
			})
		})
	})

	return result, err
}

// Prune all consumer indexes for a topic where the last timestamp is before the provided time, effectively resetting
// that consumer's index to zero.
func (bs *BoltStore) Prune(ctx context.Context, topicName string, before time.Time) ([]string, error) {
	pruned := make([]string, 0)

	err := bs.db.Update(func(tx *bbolt.Tx) error {
		topicsBucket := tx.Bucket([]byte(topicsKey))
		if topicsBucket == nil {
			return ErrNoTopic
		}

		topicBucket := topicsBucket.Bucket([]byte(topicName))
		if topicBucket == nil {
			return ErrNoTopic
		}

		consumers := topicBucket.Bucket([]byte(consumersKey))
		if consumers == nil {
			return nil
		}

		return consumers.ForEach(func(k, v []byte) error {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			var index consumer.TopicIndex
			if err := proto.Unmarshal(v, &index); err != nil {
				return fmt.Errorf("failed to unmarshal topic index data: %w", err)
			}

			if index.GetTimestamp().AsTime().After(before) {
				return nil
			}

			pruned = append(pruned, index.GetConsumerId())
			return consumers.Delete(k)
		})
	})

	return pruned, err
}
