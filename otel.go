package ae

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
)

// ErrorSpanId defines an interface for errors that can provide a span ID for distributed tracing.
type ErrorSpanId interface {
	// ErrorSpanId returns the span ID for distributed tracing.
	// Returns an empty string if no span ID is set.
	ErrorSpanId() string
}

// SpanId extracts the operation span ID from an error.
// If the error implements ErrorSpanId, returns its SpanId().
// Returns an empty string if err is nil or if the error does not implement ErrorSpanId.
func SpanId(err error) string {
	if err == nil {
		return ""
	}

	if ae, ok := err.(ErrorSpanId); ok {
		return ae.ErrorSpanId()
	}

	return ""
}

// ErrorTraceId defines an interface for errors that can provide a trace ID for distributed tracing.
type ErrorTraceId interface {
	// ErrorTraceId returns the trace ID for distributed tracing.
	// Returns an empty string if no trace ID is set.
	ErrorTraceId() string
}

// TraceId extracts the distributed tracing ID from an error.
// If the error implements ErrorTraceId, returns its TraceId().
// Returns an empty string if err is nil or if the error does not implement ErrorTraceId.
func TraceId(err error) string {
	if err == nil {
		return ""
	}

	if ae, ok := err.(ErrorTraceId); ok {
		return ae.ErrorTraceId()
	}

	return ""
}

// WithOtelAttribute returns a new context with the given OpenTelemetry attribute added.
func WithOtelAttribute(ctx context.Context, attrs ...attribute.KeyValue) context.Context {
	return WithOtelAttributes(ctx, attrs)
}

// WithOtelAttributes returns a new context with the given OpenTelemetry attributes added.
func WithOtelAttributes(ctx context.Context, attrs []attribute.KeyValue) context.Context {
	for _, attr := range attrs {
		ctx = WithAttribute(ctx, string(attr.Key), attr.Value)
	}
	return ctx
}

// WithOtelAttributeSet returns a new context with all attributes from the given OpenTelemetry attribute.Set added.
func WithOtelAttributeSet(ctx context.Context, attrs attribute.Set) context.Context {
	return WithOtelAttributes(ctx, attrs.ToSlice())
}
