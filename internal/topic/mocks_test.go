package topic_test

import (
	"context"

	"github.com/davidsbond/arrebato/internal/command"
	topicpb "github.com/davidsbond/arrebato/internal/proto/arrebato/topic/v1"
)

type (
	MockManager struct {
		err     error
		created *topicpb.Topic
		deleted string
	}

	MockExecutor struct {
		command command.Command
		err     error
	}

	MockQuerier struct {
		topic *topicpb.Topic
		err   error
	}
)

func (mm *MockQuerier) Get(ctx context.Context, name string) (*topicpb.Topic, error) {
	if mm.err != nil {
		return nil, mm.err
	}

	return mm.topic, nil
}

func (mm *MockQuerier) List(ctx context.Context) ([]*topicpb.Topic, error) {
	if mm.err != nil {
		return nil, mm.err
	}

	return []*topicpb.Topic{mm.topic}, nil
}

func (mm *MockManager) Create(ctx context.Context, t *topicpb.Topic) error {
	if mm.err != nil {
		return mm.err
	}

	mm.created = t
	return nil
}

func (mm *MockManager) Delete(ctx context.Context, t string) error {
	if mm.err != nil {
		return mm.err
	}

	mm.deleted = t
	return nil
}

func (mm *MockExecutor) Execute(ctx context.Context, cmd command.Command) error {
	if mm.err != nil {
		return mm.err
	}

	mm.command = cmd
	return nil
}
