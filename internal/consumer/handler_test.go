package consumer_test

import (
	"errors"
	"io"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"

	"github.com/davidsbond/arrebato/internal/consumer"
	consumercmd "github.com/davidsbond/arrebato/internal/proto/arrebato/consumer/command/v1"
	consumerpb "github.com/davidsbond/arrebato/internal/proto/arrebato/consumer/v1"
	"github.com/davidsbond/arrebato/internal/testutil"
)

func TestHandler_SetTopicIndex(t *testing.T) {
	t.Parallel()
	ctx := testutil.Context(t)

	tt := []struct {
		Name     string
		Command  *consumercmd.SetTopicIndex
		Expected *consumerpb.TopicIndex
		Error    error
	}{
		{
			Name: "It should create the message within the command",
			Command: &consumercmd.SetTopicIndex{
				TopicIndex: &consumerpb.TopicIndex{
					Topic: "test",
					Index: 10,
				},
			},
			Expected: &consumerpb.TopicIndex{
				Topic: "test",
				Index: 10,
			},
		},
		{
			Name: "It should return errors from the manager",
			Command: &consumercmd.SetTopicIndex{
				TopicIndex: &consumerpb.TopicIndex{
					Topic: "test",
					Index: 10,
				},
			},
			Error: io.EOF,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			mock := &MockManager{err: tc.Error}
			handler := consumer.NewHandler(mock, hclog.NewNullLogger())

			err := handler.SetTopicIndex(ctx, tc.Command)
			if tc.Error != nil {
				assert.True(t, errors.Is(err, tc.Error))
				return
			}

			assert.EqualValues(t, tc.Expected, mock.set)
		})
	}
}
