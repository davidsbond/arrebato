package tracing

import (
	"bytes"
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type (
	// The SpanFunc type is invoked when calling WithinSpan to execute some code within an OpenTelemetry span.
	SpanFunc func(ctx context.Context, span trace.Span) error
)

// WithinSpan starts a new named span and then invokes the provided SpanFunc. If the SpanFunc returns an error, it
// is added to the span and an error status is set.
func WithinSpan(ctx context.Context, name string, fn SpanFunc) error {
	ctx, span := otel.Tracer("").Start(ctx, name)
	defer span.End()

	if err := fn(ctx, span); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

// InjectBytes injects the propagation details on the span within the context and returns it as a slice of bytes.
func InjectBytes(ctx context.Context) []byte {
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	buffer := new(bytes.Buffer)
	for k, v := range carrier {
		buffer.WriteString(k)
		buffer.WriteRune(':')
		buffer.WriteString(v)
		buffer.WriteRune('\n')
	}

	return buffer.Bytes()
}

// ExtractBytes converts the span propagation details within the byte slice and returns a new context that can be
// used to create child spans of the span described within the bytes.
func ExtractBytes(ctx context.Context, b []byte) context.Context {
	m := make(map[string]string)

	lines := bytes.Split(b, []byte("\n"))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		parts := bytes.Split(line, []byte(":"))

		if len(parts) != 2 {
			return ctx
		}

		m[string(parts[0])] = string(parts[1])
	}

	return otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(m))
}

// GetSpan returns the trace.Span contained within the context.
func GetSpan(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}
