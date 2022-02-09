// Package prune provides types that remove messages from topics that have fallen outside their retention period.
package prune

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/armon/go-metrics"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"

	"github.com/davidsbond/arrebato/internal/proto/arrebato/topic/v1"
)

type (
	// The Pruner type is responsible for periodically removing messages from topics that are outside the topic's
	// configured retention period.
	Pruner struct {
		messages  MessagePruner
		topics    TopicLister
		consumers ConsumerPruner
		logger    hclog.Logger
	}

	// The TopicLister interface describes types that can list all known topics.
	TopicLister interface {
		// List should return a slice of all topics known to the current node.
		List(ctx context.Context) ([]*topic.Topic, error)
	}

	// The MessagePruner interface describes types that can remove messages from a topic.
	MessagePruner interface {
		// Prune should remove all messages from a topic older than the before time. It should return the
		// count of all pruned messages as the first return value.
		Prune(ctx context.Context, topic string, before time.Time) (uint64, error)
	}

	// The ConsumerPruner interface describes types that can reset consumer indexes for a topic. It should return
	// the names of all consumers on a topic that have been pruned as the first return value.
	ConsumerPruner interface {
		Prune(ctx context.Context, topic string, before time.Time) ([]string, error)
	}
)

// New returns a new instance of the Pruner type that will list topics via the TopicLister implementation and
// prune messages via the MessagePruner implementation.
func New(topics TopicLister, messages MessagePruner, consumers ConsumerPruner, logger hclog.Logger) *Pruner {
	return &Pruner{
		messages:  messages,
		consumers: consumers,
		topics:    topics,
		logger:    logger,
	}
}

// Prune regularly removes data from the store that has fallen outside its retention period. This method currently
// removes messages and consumer indexes that are outside a topic's configured retention period, where the retention
// period is not infinite.
func (p *Pruner) Prune(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case now := <-ticker.C:
			grp, ctx := errgroup.WithContext(ctx)
			grp.Go(func() error {
				return p.pruneMessages(ctx, now)
			})
			grp.Go(func() error {
				return p.pruneConsumers(ctx, now)
			})

			if err := grp.Wait(); err != nil {
				return fmt.Errorf("failed to prune: %w", err)
			}
		}
	}
}

func (p *Pruner) pruneMessages(ctx context.Context, now time.Time) error {
	topics, err := p.topics.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list topics: %w", err)
	}

	errs := make([]error, 0)
	for _, t := range topics {
		retentionPeriod := t.GetMessageRetentionPeriod().AsDuration()
		if retentionPeriod == 0 {
			continue
		}

		before := now.Add(-retentionPeriod)

		count, err := p.messages.Prune(ctx, t.GetName(), before)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		metrics.IncrCounterWithLabels([]string{"messages", "pruned"}, float32(count),
			[]metrics.Label{
				{
					Name:  "topic",
					Value: t.GetName(),
				},
			})
	}

	if len(errs) == 0 {
		return nil
	}

	return multierror.Append(nil, errs...).ErrorOrNil()
}

func (p *Pruner) pruneConsumers(ctx context.Context, now time.Time) error {
	topics, err := p.topics.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list topics: %w", err)
	}

	errs := make([]error, 0)
	for _, t := range topics {
		retentionPeriod := t.GetConsumerRetentionPeriod().AsDuration()
		if retentionPeriod == 0 {
			continue
		}

		before := now.Add(-retentionPeriod)

		consumers, err := p.consumers.Prune(ctx, t.GetName(), before)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		metrics.IncrCounterWithLabels([]string{"consumers", "pruned"}, float32(len(consumers)),
			[]metrics.Label{
				{
					Name:  "topic",
					Value: t.GetName(),
				},
			})
	}

	if len(errs) == 0 {
		return nil
	}

	return multierror.Append(nil, errs...).ErrorOrNil()
}
