package command_test

import (
	"io"
	"testing"
	"time"

	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/davidsbond/arrebato/internal/command"
	"github.com/davidsbond/arrebato/internal/testutil"
)

func TestExecutor_Execute(t *testing.T) {
	t.Parallel()
	ctx := testutil.Context(t)

	tt := []struct {
		Name    string
		Error   error
		Command command.Command
		State   raft.RaftState
	}{
		{
			Name:    "It should apply an encoded command if we're the raft leader",
			Command: command.New(timestamppb.New(time.Time{})),
			State:   raft.Leader,
		},
		{
			Name:    "It should return an error if we're not the leader",
			Command: command.New(timestamppb.New(time.Time{})),
			State:   raft.Follower,
			Error:   command.ErrNotLeader,
		},
		{
			Name:    "It should return other errors from the apply",
			Command: command.New(timestamppb.New(time.Time{})),
			State:   raft.Leader,
			Error:   io.EOF,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			applier := &MockApplier{err: tc.Error, state: tc.State}
			executor := command.NewExecutor(applier, time.Minute)

			err := executor.Execute(ctx, tc.Command)
			if tc.Error != nil {
				assert.EqualValues(t, tc.Error, err)
				return
			}

			actual, err := command.Unmarshal(applier.applied)
			require.NoError(t, err)
			assert.True(t, proto.Equal(tc.Command.Payload(), actual.Payload()))
		})
	}
}

func TestUnmarshal(t *testing.T) {
	expected := command.New(durationpb.New(time.Second))

	data, err := command.Marshal(expected)
	require.NoError(t, err)
	assert.NotNil(t, data)

	actual, err := command.Unmarshal(data)
	require.NoError(t, err)
	assert.True(t, proto.Equal(expected.Payload(), actual.Payload()))
}
