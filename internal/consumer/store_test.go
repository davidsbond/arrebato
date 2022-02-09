package consumer_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/davidsbond/arrebato/internal/consumer"
	consumerpb "github.com/davidsbond/arrebato/internal/proto/arrebato/consumer/v1"
	topicpb "github.com/davidsbond/arrebato/internal/proto/arrebato/topic/v1"
	"github.com/davidsbond/arrebato/internal/testutil"
	"github.com/davidsbond/arrebato/internal/topic"
)

func TestBoltStore_GetTopicIndex(t *testing.T) {
	t.Parallel()

	db := testutil.BoltDB(t)
	ctx := testutil.Context(t)

	consumers := consumer.NewBoltStore(db)
	topics := topic.NewBoltStore(db)

	const (
		topicName  = "test-topic"
		consumerID = "test-consumer"
		fixedIndex = 10
	)

	t.Run("It should return an error if no topic exists", func(t *testing.T) {
		_, err := consumers.GetTopicIndex(ctx, topicName, consumerID)

		assert.EqualValues(t, consumer.ErrNoTopic, err)
	})

	// Create a topic for future tests
	require.NoError(t, topics.Create(ctx, &topicpb.Topic{
		Name: topicName,
	}))

	t.Run("It should return zero when no consumers exist", func(t *testing.T) {
		topicIndex, err := consumers.GetTopicIndex(ctx, topicName, consumerID)

		assert.NoError(t, err)
		assert.EqualValues(t, 0, topicIndex.GetIndex())
	})

	// Set a consumer's index for future tests
	require.NoError(t, consumers.SetTopicIndex(ctx, &consumerpb.TopicIndex{
		Topic:      topicName,
		ConsumerId: consumerID,
		Index:      fixedIndex,
	}))

	t.Run("It should return the index of an existing consumer", func(t *testing.T) {
		topicIndex, err := consumers.GetTopicIndex(ctx, topicName, consumerID)

		assert.NoError(t, err)
		assert.EqualValues(t, fixedIndex, topicIndex.GetIndex())
	})

	t.Run("It should return zero for a consumer that hasn't set an index", func(t *testing.T) {
		topicIndex, err := consumers.GetTopicIndex(ctx, topicName, "new-consumer")

		assert.NoError(t, err)
		assert.Zero(t, topicIndex.GetIndex())
	})
}

func TestBoltStore_Indexes(t *testing.T) {
	t.Parallel()

	db := testutil.BoltDB(t)
	ctx := testutil.Context(t)

	consumers := consumer.NewBoltStore(db)
	topics := topic.NewBoltStore(db)

	// Set up some topics to add consumers to
	require.NoError(t, topics.Create(ctx, &topicpb.Topic{
		Name: "test-topic-1",
	}))

	require.NoError(t, topics.Create(ctx, &topicpb.Topic{
		Name: "test-topic-2",
	}))

	// Set up some consumer indexes on the two topics
	require.NoError(t, consumers.SetTopicIndex(ctx, &consumerpb.TopicIndex{
		Topic:      "test-topic-1",
		ConsumerId: "test-consumer-1",
		Index:      1000,
	}))

	require.NoError(t, consumers.SetTopicIndex(ctx, &consumerpb.TopicIndex{
		Topic:      "test-topic-1",
		ConsumerId: "test-consumer-2",
		Index:      2000,
	}))

	require.NoError(t, consumers.SetTopicIndex(ctx, &consumerpb.TopicIndex{
		Topic:      "test-topic-2",
		ConsumerId: "test-consumer-1",
		Index:      3000,
	}))

	require.NoError(t, consumers.SetTopicIndex(ctx, &consumerpb.TopicIndex{
		Topic:      "test-topic-2",
		ConsumerId: "test-consumer-2",
		Index:      4000,
	}))

	t.Run("It should return the current index of all consumers for all topics", func(t *testing.T) {
		actual, err := consumers.Indexes(ctx)
		require.NoError(t, err)

		expected := map[string]map[string]uint64{
			"test-topic-1": {
				"test-consumer-1": 1000,
				"test-consumer-2": 2000,
			},
			"test-topic-2": {
				"test-consumer-1": 3000,
				"test-consumer-2": 4000,
			},
		}

		assert.EqualValues(t, expected, actual)
	})
}

func TestBoltStore_Prune(t *testing.T) {
	t.Parallel()

	db := testutil.BoltDB(t)
	ctx := testutil.Context(t)

	consumers := consumer.NewBoltStore(db)
	topics := topic.NewBoltStore(db)

	// Set up a topic to add consumers to
	require.NoError(t, topics.Create(ctx, &topicpb.Topic{
		Name: "test-topic",
	}))

	// Set up some consumer indexes on the two topics
	require.NoError(t, consumers.SetTopicIndex(ctx, &consumerpb.TopicIndex{
		Topic:      "test-topic",
		ConsumerId: "test-consumer-1",
		Index:      1000,
		Timestamp:  timestamppb.New(time.Now().Add(time.Hour)),
	}))

	require.NoError(t, consumers.SetTopicIndex(ctx, &consumerpb.TopicIndex{
		Topic:      "test-topic",
		ConsumerId: "test-consumer-2",
		Index:      2000,
		Timestamp:  timestamppb.New(time.Now().Add(-time.Hour)),
	}))

	t.Run("It should prune consumers for a topic that haven't updated their index", func(t *testing.T) {
		pruned, err := consumers.Prune(ctx, "test-topic", time.Now())
		assert.NoError(t, err)
		assert.Contains(t, pruned, "test-consumer-2")
	})
}
