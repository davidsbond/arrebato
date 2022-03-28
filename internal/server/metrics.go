package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/armon/go-metrics"
	"github.com/armon/go-metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"

	"github.com/davidsbond/arrebato/internal/consumer"
	"github.com/davidsbond/arrebato/internal/message"
)

type (
	// The MetricConfig type contains configuration values for serving prometheus metrics.
	MetricConfig struct {
		// The Port that metrics should be served on via HTTP.
		Port int
	}

	// The MetricExporter interface describes types that export domain-specific metrics.
	MetricExporter interface {
		Export(ctx context.Context) error
	}
)

func setupMetrics() error {
	sink, err := prometheus.NewPrometheusSink()
	if err != nil {
		return fmt.Errorf("failed to initialise metrics sink: %w", err)
	}

	cnf := metrics.DefaultConfig("server")
	cnf.EnableHostname = false

	if _, err = metrics.NewGlobal(cnf, sink); err != nil {
		return fmt.Errorf("failed to configure metrics: %w", err)
	}

	return nil
}

func (svr *Server) serveMetrics(ctx context.Context, config Config) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	s := &http.Server{
		Addr:    fmt.Sprint(":", config.Metrics.Port),
		Handler: mux,
	}

	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(func() error {
		return s.ListenAndServe()
	})
	grp.Go(func() error {
		svr.raftStore.RunMetrics(ctx, time.Minute)
		return nil
	})
	grp.Go(func() error {
		<-ctx.Done()
		return s.Close()
	})

	return grp.Wait()
}

func (svr *Server) exportMetrics(ctx context.Context) error {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	exporters := []MetricExporter{
		message.NewMetricExporter(svr.messageStore),
		consumer.NewMetricExporter(svr.consumerStore),
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if !svr.IsLeader() {
				continue
			}

			for _, exporter := range exporters {
				if err := exporter.Export(ctx); err != nil {
					svr.logger.Error("failed to export metrics", "error", err.Error())
					continue
				}
			}
		}
	}
}
