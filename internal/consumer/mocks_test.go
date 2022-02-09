package consumer_test

import (
	"context"

	"github.com/davidsbond/arrebato/internal/command"
	consumerpb "github.com/davidsbond/arrebato/internal/proto/arrebato/consumer/v1"
)

type (
	MockManager struct {
		set *consumerpb.TopicIndex
		err error
	}

	MockExecutor struct {
		command command.Command
		err     error
	}
)

func (mm *MockManager) SetTopicIndex(ctx context.Context, c *consumerpb.TopicIndex) error {
	if mm.err != nil {
		return mm.err
	}

	mm.set = c
	return nil
}

func (mm *MockExecutor) Execute(ctx context.Context, cmd command.Command) error {
	if mm.err != nil {
		return mm.err
	}

	mm.command = cmd
	return nil
}
