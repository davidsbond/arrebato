package consumer

import (
	"context"
	"fmt"

	"github.com/armon/go-metrics"
)

type (
	// The MetricExporter type is used to export metrics regarding the state of consumers within the data store.
	MetricExporter struct {
		source MetricSource
	}

	// The MetricSource interface describes types that can obtain statistics about the current state of consumers that
	// can be exported as metrics.
	MetricSource interface {
		Indexes(ctx context.Context) (map[string]map[string]uint64, error)
	}
)

// NewMetricExporter returns a new instance of the MetricExporter type that will obtain statistics to export as metrics
// from the MetricSource implementation.
func NewMetricExporter(source MetricSource) *MetricExporter {
	return &MetricExporter{source: source}
}

// Export metrics regarding consumers. This method exports the current index of all consumers on all topics.
func (m *MetricExporter) Export(ctx context.Context) error {
	indexes, err := m.source.Indexes(ctx)
	if err != nil {
		return fmt.Errorf("failed to query message indexes: %w", err)
	}

	for topic, consumers := range indexes {
		for consumer, index := range consumers {
			metrics.SetGaugeWithLabels([]string{"consumer", "index"}, float32(index),
				[]metrics.Label{
					{
						Name:  "topic",
						Value: topic,
					},
					{
						Name:  "consumer_id",
						Value: consumer,
					},
				})
		}
	}

	return nil
}
