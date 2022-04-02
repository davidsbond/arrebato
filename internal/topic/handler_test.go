package topic_test

import (
	"errors"
	"io"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"

	topiccmd "github.com/davidsbond/arrebato/internal/proto/arrebato/topic/command/v1"
	topicpb "github.com/davidsbond/arrebato/internal/proto/arrebato/topic/v1"
	"github.com/davidsbond/arrebato/internal/testutil"
	"github.com/davidsbond/arrebato/internal/topic"
)

func TestHandler_Create(t *testing.T) {
	t.Parallel()
	ctx := testutil.Context(t)

	tt := []struct {
		Name     string
		Command  *topiccmd.CreateTopic
		Error    error
		Expected *topicpb.Topic
	}{
		{
			Name: "It should create a topic",
			Command: &topiccmd.CreateTopic{
				Topic: &topicpb.Topic{
					Name:                   "test-topic",
					MessageRetentionPeriod: durationpb.New(time.Minute),
				},
			},
			Expected: &topicpb.Topic{
				Name:                   "test-topic",
				MessageRetentionPeriod: durationpb.New(time.Minute),
			},
		},
		{
			Name: "It should propagate errors from the manager",
			Command: &topiccmd.CreateTopic{
				Topic: &topicpb.Topic{
					Name:                   "test-topic",
					MessageRetentionPeriod: durationpb.New(time.Minute),
				},
			},
			Error: io.EOF,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			mock := &MockStore{err: tc.Error}
			handler := topic.NewHandler(mock, hclog.NewNullLogger())

			err := handler.Create(ctx, tc.Command)
			if tc.Error != nil {
				assert.True(t, errors.Is(err, tc.Error))
				return
			}

			assert.EqualValues(t, tc.Expected, mock.created)
		})
	}
}

func TestHandler_Delete(t *testing.T) {
	t.Parallel()
	ctx := testutil.Context(t)

	tt := []struct {
		Name     string
		Command  *topiccmd.DeleteTopic
		Error    error
		Expected *topicpb.Topic
	}{
		{
			Name: "It should delete a topic",
			Command: &topiccmd.DeleteTopic{
				Name: "test-topic",
			},
			Expected: &topicpb.Topic{
				Name:                   "test-topic",
				MessageRetentionPeriod: durationpb.New(time.Minute),
			},
		},
		{
			Name: "It should propagate errors from the manager",
			Command: &topiccmd.DeleteTopic{
				Name: "test-topic",
			},
			Error: io.EOF,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			mock := &MockStore{err: tc.Error}
			handler := topic.NewHandler(mock, hclog.NewNullLogger())

			err := handler.Delete(ctx, tc.Command)
			if tc.Error != nil {
				assert.True(t, errors.Is(err, tc.Error))
				return
			}

			assert.EqualValues(t, tc.Command.GetName(), mock.deleted)
		})
	}
}
