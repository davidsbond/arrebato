package consumer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/davidsbond/arrebato/internal/command"
	"github.com/davidsbond/arrebato/internal/consumer"
	consumercmd "github.com/davidsbond/arrebato/internal/proto/arrebato/consumer/command/v1"
	consumersvc "github.com/davidsbond/arrebato/internal/proto/arrebato/consumer/service/v1"
	consumerpb "github.com/davidsbond/arrebato/internal/proto/arrebato/consumer/v1"
	"github.com/davidsbond/arrebato/internal/testutil"
)

func TestGRPC_SetTopicIndex(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name         string
		Request      *consumersvc.SetTopicIndexRequest
		ExpectedCode codes.Code
		Error        error
		Expected     command.Command
	}{
		{
			Name: "It should publish a command to set a topic index for a consumer",
			Request: &consumersvc.SetTopicIndexRequest{
				TopicIndex: &consumerpb.TopicIndex{
					Topic:      "test-topic",
					ConsumerId: "test-consumer",
					Index:      100,
				},
			},
			Expected: command.New(&consumercmd.SetTopicIndex{
				TopicIndex: &consumerpb.TopicIndex{
					Topic:      "test-topic",
					ConsumerId: "test-consumer",
					Index:      100,
				},
			}),
		},
		{
			Name: "It should return a failed precondition if the node is not the leader",
			Request: &consumersvc.SetTopicIndexRequest{
				TopicIndex: &consumerpb.TopicIndex{
					Topic:      "test-topic",
					ConsumerId: "test-consumer",
					Index:      100,
				},
			},
			ExpectedCode: codes.FailedPrecondition,
			Error:        command.ErrNotLeader,
		},
		{
			Name: "It should return not found if the topic does not exist",
			Request: &consumersvc.SetTopicIndexRequest{
				TopicIndex: &consumerpb.TopicIndex{
					Topic:      "test-topic",
					ConsumerId: "test-consumer",
					Index:      100,
				},
			},
			ExpectedCode: codes.NotFound,
			Error:        consumer.ErrNoTopic,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := testutil.Context(t)
			executor := &MockExecutor{err: tc.Error}

			resp, err := consumer.NewGRPC(executor).SetTopicIndex(ctx, tc.Request)
			require.EqualValues(t, tc.ExpectedCode, status.Code(err))

			if tc.Error != nil {
				assert.Error(t, err)
				assert.Nil(t, resp)
				return
			}

			expected := tc.Expected.Payload().(*consumercmd.SetTopicIndex)
			actual := executor.command.Payload().(*consumercmd.SetTopicIndex)

			assert.EqualValues(t, expected.GetTopicIndex().GetTopic(), actual.GetTopicIndex().GetTopic())
			assert.EqualValues(t, expected.GetTopicIndex().GetIndex(), actual.GetTopicIndex().GetIndex())
			assert.EqualValues(t, expected.GetTopicIndex().GetConsumerId(), actual.GetTopicIndex().GetConsumerId())
		})
	}
}
