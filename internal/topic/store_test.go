package topic_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"

	topicpb "github.com/davidsbond/arrebato/internal/proto/arrebato/topic/v1"
	"github.com/davidsbond/arrebato/internal/testutil"
	"github.com/davidsbond/arrebato/internal/topic"
)

func TestBoltStore_Create(t *testing.T) {
	t.Parallel()
	db := testutil.BoltDB(t)
	ctx := testutil.Context(t)
	topics := topic.NewBoltStore(db)

	const (
		topicName       = "test-topic"
		retentionPeriod = time.Hour
	)

	t.Run("It should create a topic", func(t *testing.T) {
		assert.NoError(t, topics.Create(ctx, &topicpb.Topic{
			Name:                   topicName,
			MessageRetentionPeriod: durationpb.New(retentionPeriod),
		}))
	})

	t.Run("It should return an error if the topic already exists", func(t *testing.T) {
		err := topics.Create(ctx, &topicpb.Topic{
			Name:                   topicName,
			MessageRetentionPeriod: durationpb.New(retentionPeriod),
		})

		assert.EqualValues(t, topic.ErrTopicExists, err)
	})
}

func TestBoltStore_Get(t *testing.T) {
	t.Parallel()
	db := testutil.BoltDB(t)
	ctx := testutil.Context(t)
	topics := topic.NewBoltStore(db)

	const (
		topicName       = "test-topic"
		retentionPeriod = time.Hour
	)

	t.Run("It should return a topic", func(t *testing.T) {
		expected := &topicpb.Topic{
			Name:                   topicName,
			MessageRetentionPeriod: durationpb.New(retentionPeriod),
		}

		assert.NoError(t, topics.Create(ctx, expected))
		actual, err := topics.Get(ctx, expected.GetName())
		assert.NoError(t, err)
		assert.EqualValues(t, expected.GetName(), actual.GetName())
		assert.EqualValues(t, expected.GetMessageRetentionPeriod().AsDuration(), actual.GetMessageRetentionPeriod().AsDuration())
	})

	t.Run("It should return an error if the topic does not exist", func(t *testing.T) {
		_, err := topics.Get(ctx, "nope")

		assert.EqualValues(t, topic.ErrNoTopic, err)
	})
}

func TestBoltStore_Delete(t *testing.T) {
	t.Parallel()
	db := testutil.BoltDB(t)
	ctx := testutil.Context(t)
	topics := topic.NewBoltStore(db)

	const (
		topicName       = "test-topic"
		retentionPeriod = time.Hour
	)

	t.Run("It should delete a topic", func(t *testing.T) {
		expected := &topicpb.Topic{
			Name:                   topicName,
			MessageRetentionPeriod: durationpb.New(retentionPeriod),
		}

		assert.NoError(t, topics.Create(ctx, expected))
		assert.NoError(t, topics.Delete(ctx, expected.GetName()))

		_, err := topics.Get(ctx, expected.GetName())
		assert.EqualValues(t, topic.ErrNoTopic, err)
	})

	t.Run("It should return an error if the topic does not exist", func(t *testing.T) {
		err := topics.Delete(ctx, "nope")

		assert.EqualValues(t, topic.ErrNoTopic, err)
	})
}

func TestBoltStore_List(t *testing.T) {
	t.Parallel()
	db := testutil.BoltDB(t)
	ctx := testutil.Context(t)
	topics := topic.NewBoltStore(db)

	const (
		topicName       = "test-topic"
		retentionPeriod = time.Hour
	)

	t.Run("It should list topics", func(t *testing.T) {
		expected := &topicpb.Topic{
			Name:                   topicName,
			MessageRetentionPeriod: durationpb.New(retentionPeriod),
		}

		assert.NoError(t, topics.Create(ctx, expected))

		results, err := topics.List(ctx)
		assert.NoError(t, err)
		if assert.Len(t, results, 1) {
			actual := results[0]

			assert.EqualValues(t, expected.GetName(), actual.GetName())
			assert.EqualValues(t, expected.GetMessageRetentionPeriod().AsDuration(), actual.GetMessageRetentionPeriod().AsDuration())
		}
	})
}
