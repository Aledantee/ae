package ae

import (
	"context"
	"fmt"
	"maps"
	"time"

	"go.opentelemetry.io/otel/trace"
)

// Builder is a builder for Ae errors with a fluent interface.
type Builder Ae

// New creates and returns a new instance of Builder.
func New() Builder {
	return Builder{
		tags:       make(map[string]struct{}),
		attributes: make(map[string]any),
	}
}

// Hint sets a hint message that may help resolve the error.
func (b Builder) Hint(hint string) Builder {
	b.hint = hint
	return b
}

// Timestamp sets the timestamp for when the error occurred.
func (b Builder) Timestamp(timestamp time.Time) Builder {
	b.timestamp = timestamp
	return b
}

// Now sets the current time as the error timestamp.
func (b Builder) Now() Builder {
	b.timestamp = time.Now()
	return b
}

// Code sets an error code string identifier.
func (b Builder) Code(code string) Builder {
	b.code = code
	return b
}

// ExitCode sets a non-zero exit code for the error.
// Only positive values are stored.
func (b Builder) ExitCode(exitCode int) Builder {
	if exitCode > 0 {
		b.exitCode = exitCode
	}

	return b
}

// TraceId sets the OpenTelemetry trace ID for the error.
func (b Builder) TraceId(traceId string) Builder {
	b.traceId = traceId
	return b
}

// SpanId sets the OpenTelemetry span ID for the error.
func (b Builder) SpanId(spanId string) Builder {
	b.spanId = spanId
	return b
}

// Tag adds a single tag to the error.
func (b Builder) Tag(tag string) Builder {
	b.tags[tag] = struct{}{}
	return b
}

// Tags adds multiple tags to the error.
func (b Builder) Tags(tags ...string) Builder {
	for _, tag := range tags {
		b.tags[tag] = struct{}{}
	}

	return b
}

// Attr adds a single key-value attribute to the error.
func (b Builder) Attr(key string, value any) Builder {
	b.attributes[key] = value
	return b
}

// Attrs adds multiple attributes to the error by copying from the provided map.
func (b Builder) Attrs(attrs map[string]any) Builder {
	maps.Copy(b.attributes, attrs)
	return b
}

// Cause adds one or more underlying causes to the error.
func (b Builder) Cause(causes ...error) Builder {
	return b.Causes(causes)
}

func (b Builder) Causes(causes []error) Builder {
	for _, cause := range causes {
		if cause != nil {
			b.causes = append(b.causes, cause)
		}
	}

	return b
}

// Related adds one or more related errors.
func (b Builder) Related(related ...error) Builder {
	b.related = append(b.related, related...)
	return b
}

// Stack captures the current stack trace for the error.
func (b Builder) Stack() Builder {
	b.stacks = newStack()
	return b
}

// Msg sets the error message and returns the final error.
// This is a terminal operation that completes the builder chain.
func (b Builder) Msg(msg string) error {
	b.msg = msg
	return (*Ae)(&b)
}

// UserMsg sets the error message and a user message. Then, it returns the final error.
// This is a terminal operation that completes the builder chain.
func (b Builder) UserMsg(msg, userMsg string) error {
	b.userMsg = userMsg
	return b.Msg(msg)
}

// Context extracts OpenTelemetry trace information and context values into the error.
// It captures span and trace IDs if present, and adds any requested context values as attributes.
// The keys parameter can be strings, fmt.Stringer implementations, or any other type that can be converted to a string.
func (b Builder) Context(ctx context.Context, keys ...any) Builder {
	span := trace.SpanContextFromContext(ctx)
	if span.IsValid() {
		if span.HasSpanID() {
			b.spanId = span.SpanID().String()
		}
		if span.HasTraceID() {
			b.traceId = span.TraceID().String()
		}
	}

	for _, k := range keys {
		v := ctx.Value(k)

		var (
			key   string
			value string
		)

		switch x := k.(type) {
		case string:
			key = x
		case fmt.Stringer:
			key = x.String()
		default:
			key = fmt.Sprintf("%v", x)
		}

		switch x := v.(type) {
		case string:
			value = x
		case fmt.Stringer:
			value = x.String()
		default:
			value = fmt.Sprintf("%v", x)
		}

		if key != "" && value != "" {
			b.attributes[key] = value
		}
	}

	return b
}
