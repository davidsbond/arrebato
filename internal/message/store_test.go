package message_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/davidsbond/arrebato/internal/message"
	messagepb "github.com/davidsbond/arrebato/internal/proto/arrebato/message/v1"
	topicpb "github.com/davidsbond/arrebato/internal/proto/arrebato/topic/v1"
	"github.com/davidsbond/arrebato/internal/testutil"
	"github.com/davidsbond/arrebato/internal/topic"
)

func TestBoltStore_Create(t *testing.T) {
	t.Parallel()
	db := testutil.BoltDB(t)
	ctx := testutil.Context(t)

	topics := topic.NewBoltStore(db)
	messages := message.NewBoltStore(db)

	t.Run("It should return an error if the topic does not exist", func(t *testing.T) {
		index, err := messages.Create(ctx, &messagepb.Message{
			Topic: "test-topic",
		})

		assert.EqualValues(t, message.ErrNoTopic, err)
		assert.Zero(t, index)
	})

	// Create topic for future tests
	require.NoError(t, topics.Create(ctx, &topicpb.Topic{
		Name: "test-topic",
	}))

	t.Run("It should create a message for a valid topic", func(t *testing.T) {
		index, err := messages.Create(ctx, &messagepb.Message{
			Topic:     "test-topic",
			Partition: 10,
		})

		assert.NoError(t, err)
		assert.Zero(t, index)
	})

	t.Run("It should return ordinal indexes", func(t *testing.T) {
		var lastIndex uint64

		for i := 0; i < 100; i++ {
			index, err := messages.Create(ctx, &messagepb.Message{
				Topic:     "test-topic",
				Partition: 10,
			})

			require.True(t, index == lastIndex+1)
			require.NoError(t, err)
			require.NotZero(t, index)

			lastIndex = index
		}
	})
}

func TestBoltStore_Read(t *testing.T) {
	t.Parallel()

	db := testutil.BoltDB(t)
	ctx := testutil.Context(t)

	topics := topic.NewBoltStore(db)
	messages := message.NewBoltStore(db)

	t.Run("It should return an error if the topic does not exist", func(t *testing.T) {
		err := messages.Read(ctx, "test-topic", 0, 0, nil)

		assert.EqualValues(t, message.ErrNoTopic, err)
	})

	// Create topic for future tests
	require.NoError(t, topics.Create(ctx, &topicpb.Topic{
		Name:       "test-topic",
		Partitions: 1,
	}))

	// Insert 100 test messages for us to play with
	for i := 0; i < 100; i++ {
		_, err := messages.Create(ctx, &messagepb.Message{
			Topic:     "test-topic",
			Partition: 0,
		})

		require.NoError(t, err)
	}

	t.Run("It should read messages from the start of the stream", func(t *testing.T) {
		handled := make([]*messagepb.Message, 0)

		assert.NoError(t, messages.Read(ctx, "test-topic", 0, 0, func(ctx context.Context, m *messagepb.Message) error {
			handled = append(handled, m)
			return nil
		}))

		assert.Len(t, handled, 100)
	})

	t.Run("It should read from an arbitrary point in the stream", func(t *testing.T) {
		handled := make([]*messagepb.Message, 0)

		assert.NoError(t, messages.Read(ctx, "test-topic", 0, 10, func(ctx context.Context, m *messagepb.Message) error {
			handled = append(handled, m)
			return nil
		}))

		if assert.Len(t, handled, 90) {
			assert.EqualValues(t, uint64(10), handled[0].GetIndex())
		}
	})
}

func TestBoltStore_Prune(t *testing.T) {
	t.Parallel()

	db := testutil.BoltDB(t)
	ctx := testutil.Context(t)

	topics := topic.NewBoltStore(db)
	messages := message.NewBoltStore(db)

	require.NoError(t, topics.Create(ctx, &topicpb.Topic{
		Name:       "test-topic",
		Partitions: 1,
	}))

	for i := 1; i < 11; i++ {
		ts := time.Date(2020, 1, i, 0, 0, 0, 0, time.UTC)
		_, err := messages.Create(ctx, &messagepb.Message{
			Topic:     "test-topic",
			Timestamp: timestamppb.New(ts),
			Partition: 0,
		})

		require.NoError(t, err)
	}

	t.Run("It should prune messages that were created prior to the before date", func(t *testing.T) {
		before := time.Date(2020, 1, 5, 0, 0, 0, 0, time.UTC)

		count, err := messages.Prune(ctx, "test-topic", before)
		assert.NoError(t, err)
		assert.EqualValues(t, uint64(5), count)
	})
}

func TestBoltStore_Indexes(t *testing.T) {
	t.Parallel()

	db := testutil.BoltDB(t)
	ctx := testutil.Context(t)

	topics := topic.NewBoltStore(db)
	messages := message.NewBoltStore(db)

	// Create some topics to get indexes for
	for i := 0; i < 3; i++ {
		require.NoError(t, topics.Create(ctx, &topicpb.Topic{
			Name:       strconv.Itoa(i),
			Partitions: 1,
		}))

		// Add some messages for each topic
		for j := 0; j < 10; j++ {
			_, err := messages.Create(ctx, &messagepb.Message{
				Topic:     strconv.Itoa(i),
				Partition: 0,
			})

			require.NoError(t, err)
		}
	}

	t.Run("It should return the current max index of each topic", func(t *testing.T) {
		actual, err := messages.Indexes(ctx)
		require.NoError(t, err)

		expected := map[string]map[uint32]uint64{
			"0": {
				0: 10,
			},
			"1": {
				0: 10,
			},
			"2": {
				0: 10,
			},
		}

		assert.EqualValues(t, expected, actual)
	})
}

func TestBoltStore_Counts(t *testing.T) {
	t.Parallel()

	db := testutil.BoltDB(t)
	ctx := testutil.Context(t)

	topics := topic.NewBoltStore(db)
	messages := message.NewBoltStore(db)

	// Create some topics to get indexes for
	for i := 0; i < 3; i++ {
		require.NoError(t, topics.Create(ctx, &topicpb.Topic{
			Name: strconv.Itoa(i),
		}))

		// Add some messages for each topic
		for j := 0; j < 10; j++ {
			_, err := messages.Create(ctx, &messagepb.Message{
				Topic: strconv.Itoa(i),
			})

			require.NoError(t, err)
		}
	}

	t.Run("It should return the current max index of each topic", func(t *testing.T) {
		actual, err := messages.Counts(ctx)
		require.NoError(t, err)

		expected := map[string]uint64{
			"0": 10,
			"1": 10,
			"2": 10,
		}

		assert.EqualValues(t, expected, actual)
	})
}
