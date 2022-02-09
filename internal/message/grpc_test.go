package message_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/davidsbond/arrebato/internal/command"
	"github.com/davidsbond/arrebato/internal/message"
	messagecmd "github.com/davidsbond/arrebato/internal/proto/arrebato/message/command/v1"
	messagesvc "github.com/davidsbond/arrebato/internal/proto/arrebato/message/service/v1"
	messagepb "github.com/davidsbond/arrebato/internal/proto/arrebato/message/v1"
	"github.com/davidsbond/arrebato/internal/testutil"
)

func TestGRPC_Produce(t *testing.T) {
	t.Parallel()
	ctx := testutil.Context(t)

	tt := []struct {
		Name         string
		Request      *messagesvc.ProduceRequest
		ExpectedCode codes.Code
		Error        error
		Expected     command.Command
	}{
		{
			Name: "It should execute a command to create a message",
			Request: &messagesvc.ProduceRequest{
				Message: &messagepb.Message{
					Topic:   "test-topic",
					Payload: testutil.Any(t, timestamppb.New(time.Time{})),
				},
			},
			ExpectedCode: codes.OK,
			Expected: command.New(&messagecmd.CreateMessage{
				Message: &messagepb.Message{
					Topic:   "test-topic",
					Payload: testutil.Any(t, timestamppb.New(time.Time{})),
				},
			}),
		},
		{
			Name: "It should return a failed precondition if the node is not the leader",
			Request: &messagesvc.ProduceRequest{
				Message: &messagepb.Message{
					Topic:   "test-topic",
					Payload: testutil.Any(t, timestamppb.New(time.Time{})),
				},
			},
			ExpectedCode: codes.FailedPrecondition,
			Error:        command.ErrNotLeader,
		},
		{
			Name: "It should return not found if the topic does not exist",
			Request: &messagesvc.ProduceRequest{
				Message: &messagepb.Message{
					Topic:   "test-topic",
					Payload: testutil.Any(t, timestamppb.New(time.Time{})),
				},
			},
			ExpectedCode: codes.NotFound,
			Error:        message.ErrNoTopic,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			executor := &MockExecutor{err: tc.Error}

			resp, err := message.NewGRPC(executor, nil, nil).Produce(ctx, tc.Request)
			require.EqualValues(t, tc.ExpectedCode, status.Code(err))

			if tc.Error != nil {
				assert.Error(t, err)
				assert.Nil(t, resp)
				return
			}

			expected := tc.Expected.Payload().(*messagecmd.CreateMessage)
			actual := executor.command.Payload().(*messagecmd.CreateMessage)

			assert.EqualValues(t, expected.GetMessage().GetTopic(), actual.GetMessage().GetTopic())
			assert.EqualValues(t, expected.GetMessage().GetPayload(), actual.GetMessage().GetPayload())
			assert.NotNil(t, actual.GetMessage().GetTimestamp())
		})
	}
}

func TestGRPC_Consume(t *testing.T) {
	t.Parallel()
	ctx := testutil.Context(t)

	tt := []struct {
		Name         string
		Request      *messagesvc.ConsumeRequest
		SeedTopic    []*messagepb.Message
		ExpectedCode codes.Code
		Expected     []*messagepb.Message
		Error        error
	}{
		{
			Name: "It should consume events from a valid topic",
			Request: &messagesvc.ConsumeRequest{
				Topic:      "test-topic",
				ConsumerId: "test-consumer",
			},
			SeedTopic: []*messagepb.Message{
				{
					Topic:   "test-topic",
					Payload: testutil.Any(t, timestamppb.New(time.Time{})),
				},
			},
			ExpectedCode: codes.DeadlineExceeded,
		},
		{
			Name: "It should return not found if the topic does not exist",
			Request: &messagesvc.ConsumeRequest{
				Topic:      "test-topic",
				ConsumerId: "test-consumer",
			},
			Error:        message.ErrNoTopic,
			ExpectedCode: codes.NotFound,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			// The method to stream blocks until there is a context cancellation.
			ctx, cancel := context.WithTimeout(ctx, time.Second*2)
			defer cancel()

			stream := &MockMessageConsumeServer{ctx: ctx}
			reader := &MockReader{messages: tc.SeedTopic, err: tc.Error}
			consumers := &MockTopicIndexGetter{}

			err := message.NewGRPC(nil, reader, consumers).Consume(tc.Request, stream)
			require.EqualValues(t, tc.ExpectedCode, status.Code(err))

			if tc.Error != nil {
				assert.Error(t, err)
				return
			}

			assert.EqualValues(t, tc.SeedTopic, stream.items)
		})
	}
}
