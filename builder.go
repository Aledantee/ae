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

// NewC creates and returns a new instance of Builder based on the given context.
// Shorthand for New().Context(ctx).
func NewC(ctx context.Context) Builder {
	return New().Context(ctx)
}

// From creates and returns a new instance of Builder based on the given error.
func From(err error) Builder {
	if err == nil {
		return New()
	}

	//goland:noinspection GoTypeAssertionOnErrors
	if x, ok := err.(*Ae); ok {
		return (Builder)(x.clone())
	}

	b := New()

	if x, ok := err.(ErrorMessage); ok {
		b.msg = x.ErrorMessage()
	}
	if x, ok := err.(ErrorUserMessage); ok {
		b.userMsg = x.ErrorUserMessage()
	}
	if x, ok := err.(ErrorTraceId); ok {
		b.traceId = x.ErrorTraceId()
	}
	if x, ok := err.(ErrorSpanId); ok {
		b.spanId = x.ErrorSpanId()
	}
	if x, ok := err.(ErrorTags); ok {
		b.tags = make(map[string]struct{})
		for _, tag := range x.ErrorTags() {
			b.tags[tag] = struct{}{}
		}
	}
	if x, ok := err.(ErrorCode); ok {
		b.code = x.ErrorCode()
	}
	if x, ok := err.(ErrorAttributes); ok {
		b.attributes = x.ErrorAttributes()
	}
	if x, ok := err.(ErrorExitCode); ok {
		b.exitCode = x.ErrorExitCode()
	}
	if x, ok := err.(ErrorHint); ok {
		b.hint = x.ErrorHint()
	}
	if x, ok := err.(ErrorRelated); ok {
		b.related = x.ErrorRelated()
	}
	if x, ok := err.(ErrorCauses); ok {
		b.causes = x.ErrorCauses()
	}
	if x, ok := err.(ErrorTimestamp); ok {
		b.timestamp = x.ErrorTimestamp()
	}
	if x, ok := err.(ErrorStacks); ok {
		b.stacks = x.ErrorStacks()
	}

	return b
}

// FromC creates and returns a new instance of Builder based on the given error and context.
// Shorthand for From(err).Context(ctx).
func FromC(ctx context.Context, err error) Builder {
	return From(err).Context(ctx)
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

// Causes adds one or more underlying causes to the error.
// It filters out any nil errors from the provided list.
// The causes represent errors that directly led to this error occurring.
func (b Builder) Causes(causes []error) Builder {
	for _, cause := range causes {
		if cause != nil {
			b.causes = append(b.causes, cause)
		}
	}

	return b
}

// CauseUnwrap adds one or more underlying causes to the error, unwrapping any errors that implement the Unwrap() []error interface.
// It filters out any nil errors from the provided list.
// If an error implements Unwrap() []error, its unwrapped errors are added individually.
// Otherwise, the error is added as-is.
// The causes represent errors that directly led to this error occurring.
func (b Builder) CauseUnwrap(causes ...error) Builder {
	for _, cause := range causes {
		if cause != nil {
			if x, ok := cause.(interface{ Unwrap() []error }); ok {
				for _, cause := range x.Unwrap() {
					if cause != nil {
						b.causes = append(b.causes, cause)
					}
				}
			} else {
				b.causes = append(b.causes, cause)
			}
		}
	}

	return b
}

// Related adds one or more related errors.
// It filters out any nil errors from the provided list.
// Related errors are those that are connected to this error but are not direct causes.
// This can include errors that occurred during the handling of the cause(s).
func (b Builder) Related(related ...error) Builder {
	for _, related := range related {
		if related != nil {
			b.related = append(b.related, related)
		}
	}

	return b
}

// RelatedUnwrap adds one or more related errors, unwrapping any errors that implement the Unwrap() []error interface.
// It filters out any nil errors from the provided list.
// If an error implements Unwrap() []error, its unwrapped errors are added individually.
// Otherwise, the error is added as-is.
// Related errors are those that are connected to this error but are not direct causes.
func (b Builder) RelatedUnwrap(related ...error) Builder {
	for _, related := range related {
		if related != nil {
			if x, ok := related.(interface{ Unwrap() []error }); ok {
				for _, related := range x.Unwrap() {
					if related != nil {
						b.related = append(b.related, related)
					}
				}
			} else {
				b.related = append(b.related, related)
			}
		}
	}

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

// Context extracts OpenTelemetry trace information, tags and attributes from the given context.
// Additionally, it adds the provided keys as attributes.
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

	b = b.Tags(TagsFromContext(ctx)...)
	b = b.Attrs(AttributesFromContext(ctx))

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
