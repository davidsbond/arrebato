package topic_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/davidsbond/arrebato/internal/command"
	topiccmd "github.com/davidsbond/arrebato/internal/proto/arrebato/topic/command/v1"
	topicsvc "github.com/davidsbond/arrebato/internal/proto/arrebato/topic/service/v1"
	topicpb "github.com/davidsbond/arrebato/internal/proto/arrebato/topic/v1"
	"github.com/davidsbond/arrebato/internal/testutil"
	"github.com/davidsbond/arrebato/internal/topic"
)

func TestGRPC_Create(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name         string
		Request      *topicsvc.CreateRequest
		ExpectedCode codes.Code
		Error        error
		Expected     command.Command
	}{
		{
			Name: "It should execute a command to create a topic",
			Request: &topicsvc.CreateRequest{
				Topic: &topicpb.Topic{
					Name:                   "test-topic",
					MessageRetentionPeriod: durationpb.New(time.Minute),
				},
			},
			Expected: command.New(&topiccmd.CreateTopic{
				Topic: &topicpb.Topic{
					Name:                   "test-topic",
					MessageRetentionPeriod: durationpb.New(time.Minute),
				},
			}),
		},
		{
			Name: "It should return a failed precondition if the node is not the leader",
			Request: &topicsvc.CreateRequest{
				Topic: &topicpb.Topic{
					Name:                   "test-topic",
					MessageRetentionPeriod: durationpb.New(time.Minute),
				},
			},
			ExpectedCode: codes.FailedPrecondition,
			Error:        command.ErrNotLeader,
		},
		{
			Name: "It should return already exists if the topic already exists",
			Request: &topicsvc.CreateRequest{
				Topic: &topicpb.Topic{
					Name:                   "test-topic",
					MessageRetentionPeriod: durationpb.New(time.Minute),
				},
			},
			ExpectedCode: codes.AlreadyExists,
			Error:        topic.ErrTopicExists,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := testutil.Context(t)
			executor := &MockExecutor{err: tc.Error}

			resp, err := topic.NewGRPC(executor, nil).Create(ctx, tc.Request)
			require.EqualValues(t, tc.ExpectedCode, status.Code(err))

			if tc.Error != nil {
				assert.Error(t, err)
				assert.Nil(t, resp)
				return
			}

			expected := tc.Expected.Payload().(*topiccmd.CreateTopic)
			actual := executor.command.Payload().(*topiccmd.CreateTopic)

			assert.EqualValues(t, expected.GetTopic().GetName(), actual.GetTopic().GetName())
			assert.EqualValues(t, expected.GetTopic().GetMessageRetentionPeriod(), actual.GetTopic().GetMessageRetentionPeriod())
		})
	}
}

func TestGRPC_Delete(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name         string
		Request      *topicsvc.DeleteRequest
		ExpectedCode codes.Code
		Error        error
		Expected     command.Command
	}{
		{
			Name: "It should execute a command to delete a topic",
			Request: &topicsvc.DeleteRequest{
				Name: "test-topic",
			},
			Expected: command.New(&topiccmd.DeleteTopic{
				Name: "test-topic",
			}),
		},
		{
			Name: "It should return a failed precondition if the node is not the leader",
			Request: &topicsvc.DeleteRequest{
				Name: "test-topic",
			},
			ExpectedCode: codes.FailedPrecondition,
			Error:        command.ErrNotLeader,
		},
		{
			Name: "It should return not found if the topic does not exist",
			Request: &topicsvc.DeleteRequest{
				Name: "test-topic",
			},
			ExpectedCode: codes.NotFound,
			Error:        topic.ErrNoTopic,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := testutil.Context(t)
			executor := &MockExecutor{err: tc.Error}

			resp, err := topic.NewGRPC(executor, nil).Delete(ctx, tc.Request)
			require.EqualValues(t, tc.ExpectedCode, status.Code(err))

			if tc.Error != nil {
				assert.Error(t, err)
				assert.Nil(t, resp)
				return
			}

			expected := tc.Expected.Payload().(*topiccmd.DeleteTopic)
			actual := executor.command.Payload().(*topiccmd.DeleteTopic)

			assert.EqualValues(t, expected.GetName(), actual.GetName())
		})
	}
}

func TestGRPC_Get(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name         string
		Request      *topicsvc.GetRequest
		Seed         *topicpb.Topic
		ExpectedCode codes.Code
		Error        error
		Expected     *topicsvc.GetResponse
	}{
		{
			Name: "It should return a topic",
			Request: &topicsvc.GetRequest{
				Name: "test-topic",
			},
			Seed: &topicpb.Topic{
				Name:                   "test-topic",
				MessageRetentionPeriod: durationpb.New(time.Second),
			},
			Expected: &topicsvc.GetResponse{
				Topic: &topicpb.Topic{
					Name:                   "test-topic",
					MessageRetentionPeriod: durationpb.New(time.Second),
				},
			},
		},
		{
			Name: "It should return not found if a topic does not exist",
			Request: &topicsvc.GetRequest{
				Name: "test-topic",
			},
			Error:        topic.ErrNoTopic,
			ExpectedCode: codes.NotFound,
		},
		{
			Name: "It should return failed precondition for incomplete topic info",
			Request: &topicsvc.GetRequest{
				Name: "test-topic",
			},
			Error:        topic.ErrNoTopicInfo,
			ExpectedCode: codes.FailedPrecondition,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := testutil.Context(t)

			querier := &MockQuerier{err: tc.Error, topic: tc.Seed}
			resp, err := topic.NewGRPC(nil, querier).Get(ctx, tc.Request)

			require.EqualValues(t, tc.ExpectedCode, status.Code(err))
			if tc.Error != nil {
				assert.Error(t, err)
				assert.Nil(t, resp)
				return
			}

			assert.EqualValues(t, tc.Expected, resp)
		})
	}
}

func TestGRPC_List(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name         string
		Request      *topicsvc.ListRequest
		Seed         *topicpb.Topic
		ExpectedCode codes.Code
		Error        error
		Expected     *topicsvc.ListResponse
	}{
		{
			Name:    "It should return topics",
			Request: &topicsvc.ListRequest{},
			Seed: &topicpb.Topic{
				Name:                   "test-topic",
				MessageRetentionPeriod: durationpb.New(time.Second),
			},
			Expected: &topicsvc.ListResponse{
				Topics: []*topicpb.Topic{
					{
						Name:                   "test-topic",
						MessageRetentionPeriod: durationpb.New(time.Second),
					},
				},
			},
		},
		{
			Name:         "It should return failed precondition for incomplete topic info",
			Request:      &topicsvc.ListRequest{},
			Error:        topic.ErrNoTopicInfo,
			ExpectedCode: codes.FailedPrecondition,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := testutil.Context(t)

			querier := &MockQuerier{err: tc.Error, topic: tc.Seed}
			resp, err := topic.NewGRPC(nil, querier).List(ctx, tc.Request)

			require.EqualValues(t, tc.ExpectedCode, status.Code(err))
			if tc.Error != nil {
				assert.Error(t, err)
				assert.Nil(t, resp)
				return
			}

			assert.EqualValues(t, tc.Expected, resp)
		})
	}
}
