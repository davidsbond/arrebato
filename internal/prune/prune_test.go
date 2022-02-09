package prune_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/davidsbond/arrebato/internal/consumer"
	"github.com/davidsbond/arrebato/internal/message"
	consumerpb "github.com/davidsbond/arrebato/internal/proto/arrebato/consumer/v1"
	messagepb "github.com/davidsbond/arrebato/internal/proto/arrebato/message/v1"
	topicpb "github.com/davidsbond/arrebato/internal/proto/arrebato/topic/v1"
	"github.com/davidsbond/arrebato/internal/prune"
	"github.com/davidsbond/arrebato/internal/testutil"
	"github.com/davidsbond/arrebato/internal/topic"
)

func TestPruner_PruneMessages(t *testing.T) {
	t.Parallel()

	db := testutil.BoltDB(t)
	ctx := testutil.Context(t)

	topics := topic.NewBoltStore(db)

	// Create a topic to use in tests
	require.NoError(t, topics.Create(ctx, &topicpb.Topic{
		Name:                   "test-topic",
		MessageRetentionPeriod: durationpb.New(time.Second),
	}))

	messages := message.NewBoltStore(db)
	consumers := consumer.NewBoltStore(db)

	// Put a message on the topic that we expect to be pruned
	_, err := messages.Create(ctx, &messagepb.Message{
		Topic:     "test-topic",
		Payload:   testutil.Any(t, structpb.NewStringValue("hello")),
		Timestamp: timestamppb.New(time.Now().Add(-time.Hour)),
	})
	require.NoError(t, err)

	// Put a message on the topic we expect to still be there
	_, err = messages.Create(ctx, &messagepb.Message{
		Topic:     "test-topic",
		Payload:   testutil.Any(t, structpb.NewStringValue("world")),
		Timestamp: timestamppb.New(time.Now().Add(time.Hour)),
	})
	require.NoError(t, err)

	pruner := prune.New(topics, messages, consumers, hclog.NewNullLogger())

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		_ = pruner.Prune(ctx, time.Second)
	}()

	// Keep reading messages, eventually the pruning will kick in and the first message we get on a call to Read will
	// be the second one we produced with the string "world" as the payload.
	assert.Eventually(t, func() bool {
		var gotMessage bool
		_ = messages.Read(ctx, "test-topic", 0, func(ctx context.Context, m *messagepb.Message) error {
			payload, err := m.GetPayload().UnmarshalNew()
			require.NoError(t, err)

			switch p := payload.(type) {
			case *structpb.Value:
				gotMessage = p.GetStringValue() == "world"
			default:
				t.Fatalf("expected *structpb.Value, got %T", payload)
			}

			// Return an arbitrary error here so that we only check the first message we're given.
			return io.EOF
		})

		return gotMessage
	}, time.Minute, time.Second)
}

func TestPruner_PruneConsumerIndexes(t *testing.T) {
	t.Parallel()

	db := testutil.BoltDB(t)
	ctx := testutil.Context(t)

	topics := topic.NewBoltStore(db)

	// Create a topic to use in tests
	require.NoError(t, topics.Create(ctx, &topicpb.Topic{
		Name:                    "test-topic",
		ConsumerRetentionPeriod: durationpb.New(time.Second),
	}))

	messages := message.NewBoltStore(db)
	consumers := consumer.NewBoltStore(db)

	// Set a consumer index we expect to be reset after pruning
	require.NoError(t, consumers.SetTopicIndex(ctx, &consumerpb.TopicIndex{
		Topic:      "test-topic",
		ConsumerId: "test-consumer",
		Index:      100,
		Timestamp:  timestamppb.Now(),
	}))

	pruner := prune.New(topics, messages, consumers, hclog.NewNullLogger())

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		_ = pruner.Prune(ctx, time.Second)
	}()

	// Eventually the index for this consumer on the topic should be reset to zero.
	assert.Eventually(t, func() bool {
		idx, err := consumers.GetTopicIndex(ctx, "test-topic", "test-consumer")
		assert.NoError(t, err)
		return idx.GetIndex() == 0
	}, time.Minute, time.Second)
}
