package message_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/davidsbond/arrebato/internal/command"
	"github.com/davidsbond/arrebato/internal/message"
	messagecmd "github.com/davidsbond/arrebato/internal/proto/arrebato/message/command/v1"
	messagesvc "github.com/davidsbond/arrebato/internal/proto/arrebato/message/service/v1"
	messagepb "github.com/davidsbond/arrebato/internal/proto/arrebato/message/v1"
	topicpb "github.com/davidsbond/arrebato/internal/proto/arrebato/topic/v1"
	"github.com/davidsbond/arrebato/internal/testutil"
)

func TestGRPC_Produce(t *testing.T) {
	t.Parallel()
	ctx := testutil.Context(t)

	tt := []struct {
		Name            string
		Request         *messagesvc.ProduceRequest
		RequireVerified bool
		ACLAllow        bool
		ExpectedCode    codes.Code
		Error           error
		Expected        command.Command
	}{
		{
			Name:     "It should execute a command to create a message",
			ACLAllow: true,
			Request: &messagesvc.ProduceRequest{
				Message: &messagepb.Message{
					Topic: "test-topic",
					Value: testutil.Any(t, timestamppb.New(time.Time{})),
				},
			},
			ExpectedCode: codes.OK,
			Expected: command.New(&messagecmd.CreateMessage{
				Message: &messagepb.Message{
					Topic:     "test-topic",
					Value:     testutil.Any(t, timestamppb.New(time.Time{})),
					Partition: 1,
				},
			}),
		},
		{
			Name:     "It should return a failed precondition if the node is not the leader",
			ACLAllow: true,
			Request: &messagesvc.ProduceRequest{
				Message: &messagepb.Message{
					Topic: "test-topic",
					Value: testutil.Any(t, timestamppb.New(time.Time{})),
				},
			},
			ExpectedCode: codes.FailedPrecondition,
			Error:        command.ErrNotLeader,
		},
		{
			Name:     "It should return not found if the topic does not exist",
			ACLAllow: true,
			Request: &messagesvc.ProduceRequest{
				Message: &messagepb.Message{
					Topic: "test-topic",
					Value: testutil.Any(t, timestamppb.New(time.Time{})),
				},
			},
			ExpectedCode: codes.NotFound,
			Error:        message.ErrNoTopic,
		},
		{
			Name: "It should return permission denied if the ACL does not allow producing",
			Request: &messagesvc.ProduceRequest{
				Message: &messagepb.Message{
					Topic: "test-topic",
					Value: testutil.Any(t, timestamppb.New(time.Time{})),
				},
			},
			ExpectedCode: codes.PermissionDenied,
		},
		{
			Name:            "It should return permission denied if the message is not verified and the topic requires it",
			RequireVerified: true,
			ACLAllow:        true,
			Request: &messagesvc.ProduceRequest{
				Message: &messagepb.Message{
					Topic: "test-topic",
					Value: testutil.Any(t, timestamppb.New(time.Time{})),
					Key:   testutil.Any(t, timestamppb.New(time.Time{})),
				},
			},
			ExpectedCode: codes.PermissionDenied,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			executor := &MockExecutor{err: tc.Error}
			publicKeys := &MockPublicKeyGetter{}
			topics := &MockTopicGetter{topic: &topicpb.Topic{RequireVerifiedMessages: tc.RequireVerified, Partitions: 10}}

			svc := message.NewGRPC(
				executor,
				nil,
				nil,
				&MockACL{allowed: tc.ACLAllow},
				publicKeys,
				topics,
				message.NewCRC32Partitioner(),
				"",
				nil,
			)

			resp, err := svc.Produce(ctx, tc.Request)
			require.EqualValues(t, tc.ExpectedCode, status.Code(err))

			if tc.Error != nil || tc.ExpectedCode > codes.OK {
				assert.Error(t, err)
				assert.Nil(t, resp)
				return
			}

			expected := tc.Expected.Payload().(*messagecmd.CreateMessage)
			actual := executor.command.Payload().(*messagecmd.CreateMessage)

			assert.EqualValues(t, expected.GetMessage().GetTopic(), actual.GetMessage().GetTopic())
			assert.EqualValues(t, expected.GetMessage().GetPartition(), actual.GetMessage().GetPartition())
			assert.True(t, proto.Equal(expected.GetMessage().GetValue(), actual.GetMessage().GetValue()))
			assert.True(t, proto.Equal(expected.GetMessage().GetKey(), actual.GetMessage().GetKey()))
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
		ACLAllow     bool
		Allocated    bool
		SeedTopic    []*messagepb.Message
		ExpectedCode codes.Code
		Expected     []*messagepb.Message
		Error        error
	}{
		{
			Name:      "It should consume events from a valid topic",
			ACLAllow:  true,
			Allocated: true,
			Request: &messagesvc.ConsumeRequest{
				Topic:      "test-topic",
				ConsumerId: "test-consumer",
			},
			SeedTopic: []*messagepb.Message{
				{
					Topic: "test-topic",
					Value: testutil.Any(t, timestamppb.New(time.Time{})),
				},
			},
			ExpectedCode: codes.DeadlineExceeded,
		},
		{
			Name:      "It should return not found if the topic does not exist",
			ACLAllow:  true,
			Allocated: true,
			Request: &messagesvc.ConsumeRequest{
				Topic:      "test-topic",
				ConsumerId: "test-consumer",
			},
			Error:        message.ErrNoTopic,
			ExpectedCode: codes.NotFound,
		},
		{
			Name:      "It should return permission denied if the ACL does not allow consuming",
			Allocated: true,
			Request: &messagesvc.ConsumeRequest{
				Topic:      "test-topic",
				ConsumerId: "test-consumer",
			},
			ExpectedCode: codes.PermissionDenied,
		},
		{
			Name:     "It should return failed precondition if the topic is not allocated to the node",
			ACLAllow: true,
			Request: &messagesvc.ConsumeRequest{
				Topic:      "test-topic",
				ConsumerId: "test-consumer",
			},
			ExpectedCode: codes.FailedPrecondition,
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
			nodes := &MockNodeStore{allocated: tc.Allocated}

			err := message.NewGRPC(
				nil,
				reader,
				consumers,
				&MockACL{allowed: tc.ACLAllow},
				nil,
				nil,
				nil,
				"",
				nodes).Consume(tc.Request, stream)

			require.EqualValues(t, tc.ExpectedCode, status.Code(err))

			if tc.Error != nil || tc.ExpectedCode > codes.OK {
				assert.Error(t, err)
				return
			}

			assert.EqualValues(t, tc.SeedTopic, stream.items)
		})
	}
}
