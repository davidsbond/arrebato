package server

import (
	"context"
	"fmt"
	"io"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

type (
	tracerCloser struct {
		provider *sdktrace.TracerProvider
	}

	// The TracingConfig type contains configuration values for OpenTelemetry tracing to a collector.
	TracingConfig struct {
		// The HTTP(s) endpoint to send OpenTelemetry spans to.
		Endpoint string
		// Whether tracing is enabled.
		Enabled bool
	}
)

func (t *tracerCloser) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return t.provider.Shutdown(ctx)
}

func setupTracing(ctx context.Context, config Config) (io.Closer, error) {
	if !config.Tracing.Enabled {
		return io.NopCloser(nil), nil
	}

	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(config.Tracing.Endpoint),
		otlptracehttp.WithInsecure(),
	)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.AdvertiseAddress),
			semconv.ServiceVersionKey.String(config.Version),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tracing resource: %w", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return &tracerCloser{provider: provider}, nil
}
