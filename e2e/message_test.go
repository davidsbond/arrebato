package e2e_test

import (
	"context"
	"strconv"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/davidsbond/arrebato/internal/testutil"
	"github.com/davidsbond/arrebato/pkg/arrebato"
)

type (
	MessageSuite struct {
		suite.Suite

		client *arrebato.Client
	}
)

func (s *MessageSuite) TestProduceConsumeMessages() {
	ctx := testutil.Context(s.T())

	// Create a topic for use in the test suite
	require.NoError(s.T(), s.client.CreateTopic(ctx, arrebato.Topic{
		Name:                   "test-suite-topic",
		MessageRetentionPeriod: time.Hour,
	}))

	expected := structpb.NewStringValue("hello-world")

	s.Run("It should produce a message on a topic", func() {
		producer := s.client.NewProducer("test-suite-topic")

		require.NoError(s.T(), producer.Produce(ctx, arrebato.Message{
			Payload: expected,
		}))
	})

	s.Run("It should consume the produced message from the topic", func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*2)
		defer cancel()

		consumer, err := s.client.NewConsumer(ctx, arrebato.ConsumerConfig{
			Topic:      "test-suite-topic",
			ConsumerID: "test-suite-consumer",
		})
		require.NoError(s.T(), err)
		s.T().Cleanup(func() {
			assert.NoError(s.T(), consumer.Close())
		})

		err = consumer.Consume(ctx, func(ctx context.Context, m proto.Message) error {
			switch actual := m.(type) {
			case *structpb.Value:
				assert.EqualValues(s.T(), expected.GetStringValue(), actual.GetStringValue())
				return nil
			default:
				s.Failf("unexpected type", "expected *structpb.Value, got %T", m)
				return nil
			}
		})

		assert.Error(s.T(), err)
	})

	s.Run("It should produce many messages on a topic", func() {
		producer := s.client.NewProducer("test-suite-topic")

		for i := 0; i < 1024; i++ {
			expected = structpb.NewStringValue(strconv.Itoa(i))

			require.NoError(s.T(), producer.Produce(ctx, arrebato.Message{
				Payload: expected,
			}))
		}
	})

	s.Run("It should resume from its last place in the stream", func() {
		ctx, cancel := context.WithCancel(ctx)

		consumer, err := s.client.NewConsumer(ctx, arrebato.ConsumerConfig{
			Topic:      "test-suite-topic",
			ConsumerID: "test-suite-consumer",
		})
		require.NoError(s.T(), err)
		s.T().Cleanup(func() {
			assert.NoError(s.T(), consumer.Close())
		})

		err = consumer.Consume(ctx, func(ctx context.Context, m proto.Message) error {
			defer cancel()

			switch actual := m.(type) {
			case *structpb.Value:
				// We should have the first message from the 1024 we produced earlier, not the "hello-world" one we
				// initially produced.
				assert.EqualValues(s.T(), "0", actual.GetStringValue())
				return nil
			default:
				s.Failf("unexpected type", "expected *structpb.Value, got %T", m)
				return nil
			}
		})

		assert.Error(s.T(), err)
	})
}
