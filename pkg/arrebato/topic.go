package arrebato

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"

	topicsvc "github.com/davidsbond/arrebato/internal/proto/arrebato/topic/service/v1"
	"github.com/davidsbond/arrebato/internal/proto/arrebato/topic/v1"
)

type (
	// The Topic type describes a topic stored within the cluster.
	Topic struct {
		// The Name of the Topic.
		Name string

		// The amount of time messages on the Topic will be stored.
		MessageRetentionPeriod time.Duration

		// The maximum age of a consumer index on a Topic before it is reset to zero.
		ConsumerRetentionPeriod time.Duration
	}
)

var (
	// ErrNoTopic is the error given when attempting to perform an operation against a topic that does not exist.
	ErrNoTopic = errors.New("no topic")

	// ErrTopicExists is the error given when attempting to create a topic that already exists.
	ErrTopicExists = errors.New("topic exists")
)

// CreateTopic creates a new topic described by the provided Topic. Returns ErrTopicExists if the topic already
// exists.
func (c *Client) CreateTopic(ctx context.Context, t Topic) error {
	svc := topicsvc.NewTopicServiceClient(c.cluster.leader())
	_, err := svc.Create(ctx, &topicsvc.CreateRequest{
		Topic: &topic.Topic{
			Name:                    t.Name,
			MessageRetentionPeriod:  durationpb.New(t.MessageRetentionPeriod),
			ConsumerRetentionPeriod: durationpb.New(t.ConsumerRetentionPeriod),
		},
	})
	switch {
	case status.Code(err) == codes.AlreadyExists:
		return ErrTopicExists
	case status.Code(err) == codes.FailedPrecondition:
		c.cluster.findLeader(ctx)
		return c.CreateTopic(ctx, t)
	case err != nil:
		return err
	default:
		return nil
	}
}

// Topic returns a named Topic. Returns ErrNoTopic if the topic does not exist.
func (c *Client) Topic(ctx context.Context, name string) (Topic, error) {
	svc := topicsvc.NewTopicServiceClient(c.cluster.any())
	resp, err := svc.Get(ctx, &topicsvc.GetRequest{
		Name: name,
	})
	switch {
	case status.Code(err) == codes.NotFound:
		return Topic{}, ErrNoTopic
	case err != nil:
		return Topic{}, err
	default:
		return Topic{
			Name:                    resp.GetTopic().GetName(),
			MessageRetentionPeriod:  resp.GetTopic().GetMessageRetentionPeriod().AsDuration(),
			ConsumerRetentionPeriod: resp.GetTopic().GetConsumerRetentionPeriod().AsDuration(),
		}, nil
	}
}

// Topics lists all topics stored in the cluster.
func (c *Client) Topics(ctx context.Context) ([]Topic, error) {
	svc := topicsvc.NewTopicServiceClient(c.cluster.any())
	resp, err := svc.List(ctx, &topicsvc.ListRequest{})
	if err != nil {
		return nil, err
	}

	topics := make([]Topic, len(resp.GetTopics()))
	for i, tp := range resp.GetTopics() {
		topics[i] = Topic{
			Name:                    tp.GetName(),
			MessageRetentionPeriod:  tp.GetMessageRetentionPeriod().AsDuration(),
			ConsumerRetentionPeriod: tp.GetConsumerRetentionPeriod().AsDuration(),
		}
	}

	return topics, nil
}

// DeleteTopic removes a named Topic from the cluster. Returns ErrNoTopic if the topic does not exist.
func (c *Client) DeleteTopic(ctx context.Context, name string) error {
	svc := topicsvc.NewTopicServiceClient(c.cluster.leader())
	_, err := svc.Delete(ctx, &topicsvc.DeleteRequest{
		Name: name,
	})

	switch {
	case status.Code(err) == codes.NotFound:
		return ErrNoTopic
	case status.Code(err) == codes.FailedPrecondition:
		c.cluster.findLeader(ctx)
		return c.DeleteTopic(ctx, name)
	case err != nil:
		return err
	default:
		return nil
	}
}
