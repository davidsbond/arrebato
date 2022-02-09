package message_test

import (
	"errors"
	"io"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"

	"github.com/davidsbond/arrebato/internal/message"
	messagecmd "github.com/davidsbond/arrebato/internal/proto/arrebato/message/command/v1"
	messagepb "github.com/davidsbond/arrebato/internal/proto/arrebato/message/v1"
	"github.com/davidsbond/arrebato/internal/testutil"
)

func TestHandler_Create(t *testing.T) {
	t.Parallel()
	ctx := testutil.Context(t)

	tt := []struct {
		Name     string
		Command  *messagecmd.CreateMessage
		Expected *messagepb.Message
		Error    error
	}{
		{
			Name: "It should create the message within the command",
			Command: &messagecmd.CreateMessage{
				Message: &messagepb.Message{
					Topic: "test",
				},
			},
			Expected: &messagepb.Message{
				Topic: "test",
			},
		},
		{
			Name: "It should return errors from the manager",
			Command: &messagecmd.CreateMessage{
				Message: &messagepb.Message{
					Topic: "test",
				},
			},
			Error: io.EOF,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			mock := &MockCreator{err: tc.Error}
			handler := message.NewHandler(mock, hclog.NewNullLogger())

			err := handler.Create(ctx, tc.Command)
			if tc.Error != nil {
				assert.True(t, errors.Is(err, tc.Error))
				return
			}

			assert.EqualValues(t, tc.Expected, mock.created)
		})
	}
}
