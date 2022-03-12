package e2e_test

import (
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/davidsbond/arrebato/internal/testutil"
	"github.com/davidsbond/arrebato/pkg/arrebato"
)

type (
	TopicSuite struct {
		suite.Suite

		client *arrebato.Client
	}
)

func (s *TopicSuite) TestManageTopics() {
	ctx := testutil.Context(s.T())

	s.Run("It should create a new topic", func() {
		topic := arrebato.Topic{
			Name:                   "test-topic",
			MessageRetentionPeriod: time.Hour,
		}

		require.NoError(s.T(), s.client.CreateTopic(ctx, topic))
	})

	s.Run("It should return topic information", func() {
		topic, err := s.client.Topic(ctx, "test-topic")

		assert.NoError(s.T(), err)
		assert.EqualValues(s.T(), "test-topic", topic.Name)
		assert.EqualValues(s.T(), time.Hour, topic.MessageRetentionPeriod)
	})

	s.Run("It should return an error requesting information on a topic that does not exist", func() {
		_, err := s.client.Topic(ctx, "non-existent")
		assert.Error(s.T(), err)
	})

	s.Run("It should list all topics", func() {
		topics, err := s.client.Topics(ctx)

		assert.NoError(s.T(), err)
		assert.Len(s.T(), topics, 1)
	})

	s.Run("It should delete a topic", func() {
		assert.NoError(s.T(), s.client.DeleteTopic(ctx, "test-topic"))
	})

	s.Run("It should return an error deleting a topic that does not exist", func() {
		assert.Error(s.T(), s.client.DeleteTopic(ctx, "non-existent"))
	})
}
