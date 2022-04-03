package message_test

import (
	"context"

	"github.com/davidsbond/arrebato/internal/command"
	"github.com/davidsbond/arrebato/internal/message"
	"github.com/davidsbond/arrebato/internal/proto/arrebato/acl/v1"
	consumerpb "github.com/davidsbond/arrebato/internal/proto/arrebato/consumer/v1"
	messagesvc "github.com/davidsbond/arrebato/internal/proto/arrebato/message/service/v1"
	messabepb "github.com/davidsbond/arrebato/internal/proto/arrebato/message/v1"
	topicpb "github.com/davidsbond/arrebato/internal/proto/arrebato/topic/v1"
)

type (
	MockCreator struct {
		created *messabepb.Message
		err     error
	}

	MockExecutor struct {
		command command.Command
		err     error
	}

	MockMessageConsumeServer struct {
		messagesvc.MessageService_ConsumeServer

		items []*messabepb.Message
		err   error
		ctx   context.Context
	}

	MockReader struct {
		messages []*messabepb.Message
		err      error
	}

	MockTopicIndexGetter struct{}

	MockACL struct {
		allowed bool
	}

	MockPublicKeyGetter struct {
		err error
		key []byte
	}

	MockTopicGetter struct {
		topic *topicpb.Topic
	}
)

func (mm *MockTopicGetter) Get(ctx context.Context, name string) (*topicpb.Topic, error) {
	return mm.topic, nil
}

func (mm *MockPublicKeyGetter) Get(ctx context.Context, clientID string) ([]byte, error) {
	return mm.key, mm.err
}

func (mm *MockACL) Allowed(ctx context.Context, topic, client string, permission acl.Permission) (bool, error) {
	return mm.allowed, nil
}

func (mm *MockTopicIndexGetter) GetTopicIndex(ctx context.Context, topic string, consumerID string) (*consumerpb.TopicIndex, error) {
	return &consumerpb.TopicIndex{
		Topic:      topic,
		ConsumerId: consumerID,
		Index:      0,
	}, nil
}

func (mm *MockReader) Read(ctx context.Context, topic string, partition uint32, startIndex uint64, fn message.ReadFunc) error {
	if mm.err != nil {
		return mm.err
	}

	for _, m := range mm.messages {
		if m.Index < startIndex {
			continue
		}

		if err := fn(ctx, m); err != nil {
			return err
		}
	}

	return nil
}

func (mm *MockMessageConsumeServer) Send(response *messagesvc.ConsumeResponse) error {
	if mm.err != nil {
		return mm.err
	}

	mm.items = append(mm.items, response.GetMessage())
	return nil
}

func (mm *MockMessageConsumeServer) Context() context.Context {
	return mm.ctx
}

func (mm *MockExecutor) Execute(ctx context.Context, cmd command.Command) error {
	if mm.err != nil {
		return mm.err
	}

	mm.command = cmd
	return nil
}

func (mm *MockCreator) Create(ctx context.Context, m *messabepb.Message) (uint64, error) {
	if mm.err != nil {
		return 0, mm.err
	}

	mm.created = m
	return 1, nil
}
