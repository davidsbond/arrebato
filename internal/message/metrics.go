package message

import (
	"context"
	"fmt"

	"github.com/armon/go-metrics"
)

type (
	// The MetricExporter type is used to export metrics regarding the state of messages within the data store.
	MetricExporter struct {
		source MetricSource
	}

	// The MetricSource interface describes types that can obtain statistics about the current state of messages that
	// can be exported as metrics.
	MetricSource interface {
		Indexes(ctx context.Context) (map[string]uint64, error)
		Counts(ctx context.Context) (map[string]uint64, error)
	}
)

// NewMetricExporter returns a new instance of the MetricExporter type that will obtain statistics to export as metrics
// from the MetricSource implementation.
func NewMetricExporter(source MetricSource) *MetricExporter {
	return &MetricExporter{source: source}
}

// Export metrics regarding messages. This method exports the current index of all topics as well as the count of
// messages within those topics. Since topics can be pruned, the current index does not always reflect the total number
// of messages that can be consumed.
func (m *MetricExporter) Export(ctx context.Context) error {
	indexes, err := m.source.Indexes(ctx)
	if err != nil {
		return fmt.Errorf("failed to query message indexes: %w", err)
	}

	for topic, index := range indexes {
		metrics.SetGaugeWithLabels([]string{"message", "index"}, float32(index),
			[]metrics.Label{
				{
					Name:  "topic",
					Value: topic,
				},
			})
	}

	counts, err := m.source.Counts(ctx)
	if err != nil {
		return fmt.Errorf("failed to query message counts: %w", err)
	}

	for topic, count := range counts {
		metrics.SetGaugeWithLabels([]string{"message", "count"}, float32(count),
			[]metrics.Label{
				{
					Name:  "topic",
					Value: topic,
				},
			})
	}

	return nil
}
