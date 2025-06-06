package ae

import (
	"context"
	"fmt"
	"syscall"

	"go.opentelemetry.io/otel/trace"
)

// ErrorBuilder is a builder type for constructing Error instances.
// It provides a fluent interface for setting various error properties.
type ErrorBuilder Error

// New creates a new error builder with the given message.
// The message must not be empty. If an empty message is provided, New panics.
func New(msg string) *ErrorBuilder {
	if msg == "" {
		panic("ae: error message must not be empty")
	}

	return newInternal(msg)
}

// Newf creates a new error builder with a formatted message.
// It uses fmt.Sprintf to format the message with the provided arguments.
// The formatted message must not be empty. If the resulting message is empty, Newf panics.
func Newf(format string, args ...any) *ErrorBuilder {
	return New(fmt.Sprintf(format, args...))
}

// NotImplemeted creates a new error builder with the message "not implemented".
// This is a common error message for when a feature is not yet implemented.
func NotImplemeted() *ErrorBuilder {
	return New("not implemented")
}

// newInternal creates a new ErrorBuilder with default values.
// It initializes maps for tags and attributes, and sets a default exit code of 1.
func newInternal(msg string) *ErrorBuilder {
	return &ErrorBuilder{
		msg:      msg,
		tags:     make(map[string]struct{}),
		attrs:    make(map[string]any),
		exitCode: 1,
	}
}

// Msg sets the internal error message.
// If the provided message is empty, the current message is preserved.
func (e *ErrorBuilder) Msg(msg string) *ErrorBuilder {
	if msg != "" {
		e.msg = msg
	}

	return e
}

// Public sets the public-facing error message.
// This message is suitable for end users and should not contain sensitive information.
func (e *ErrorBuilder) Public(msg string) *ErrorBuilder {
	e.pubMsg = msg
	return e
}

// Hint sets a hint message that provides guidance for resolving the error.
func (e *ErrorBuilder) Hint(hint string) *ErrorBuilder {
	e.hint = hint
	return e
}

// Code sets an error code that can be used for programmatic error handling.
func (e *ErrorBuilder) Code(code string) *ErrorBuilder {
	e.code = code
	return e
}

// ExitCode sets the process exit code associated with this error.
// Only positive values are accepted; negative values are ignored.
func (e *ErrorBuilder) ExitCode(code int) *ErrorBuilder {
	if code > 0 {
		e.exitCode = code
	}

	return e
}

func (e *ErrorBuilder) TraceId(traceId string) *ErrorBuilder {
	if e.traceId != "" {
		e.traceId = traceId
	}
	return e
}

func (e *ErrorBuilder) SpanId(spanId string) *ErrorBuilder {
	if e.spanId != "" {
		e.spanId = spanId
	}
	return e
}

// Tag adds a tag to categorize or filter the error.
func (e *ErrorBuilder) Tag(tags ...string) *ErrorBuilder {
	for _, t := range tags {
		e.tags[t] = struct{}{}
	}
	return e
}

// Attr sets a single attribute with the given key and value.
func (e *ErrorBuilder) Attr(key string, value any) *ErrorBuilder {
	if key != "" {
		e.attrs[key] = value
	}
	return e
}

// Attrs sets multiple attributes from key-value pairs.
// If an odd number of arguments is provided, the last value defaults to "!VALUE!".
// Keys are converted to strings using fmt.Stringer if available, otherwise using fmt.Sprintf.
func (e *ErrorBuilder) Attrs(kv ...any) *ErrorBuilder {
	if len(kv)%2 != 0 {
		kv = append(kv, "!VALUE!")
	}

	for i := 0; i < len(kv); i += 2 {
		var key string

		switch x := kv[i].(type) {
		case string:
			key = x
		case fmt.Stringer:
			key = x.String()
		default:
			key = fmt.Sprintf("%v", x)
		}

		if key != "" {
			e.attrs[key] = kv[i+1]
		}
	}

	return e
}

func (e *ErrorBuilder) AttrsMap(m map[string]any) *ErrorBuilder {
	for k, v := range m {
		if k != "" {
			e.attrs[k] = v
		}
	}
	return e
}

// Cause adds one or more errors as direct causes of this error.
func (e *ErrorBuilder) Cause(errs ...error) *ErrorBuilder {
	for _, err := range errs {
		if err != nil {
			e.causes = append(e.causes, err)
		}
	}
	return e
}

// Related adds one or more errors that are related to this error.
func (e *ErrorBuilder) Related(errs ...error) *ErrorBuilder {
	for _, err := range errs {
		if err != nil {
			e.relatedErrs = append(e.relatedErrs, err)
		}
	}
	return e
}

// Context adds tracing information and context values from the provided context.
// It extracts the trace ID and span ID if available in the context.
// For each key provided, it adds the corresponding value from the context as an attribute.
// Keys can be strings, fmt.Stringer implementations, or any other type that can be converted to a string. Empty keys are ignored.
func (e *ErrorBuilder) Context(ctx context.Context, keys ...any) *ErrorBuilder {
	sp := trace.SpanContextFromContext(ctx)
	if sp.IsValid() {
		e.traceId = sp.TraceID().String()
		e.spanId = sp.SpanID().String()
	}

	for _, k := range keys {
		var key string

		switch x := k.(type) {
		case string:
			key = x
		case fmt.Stringer:
			key = x.String()
		default:
			key = fmt.Sprintf("%v", x)
		}

		if key != "" {
			e.Attr(key, ctx.Value(k))
		}
	}

	return e
}

// Build converts the ErrorBuilder to an Error.
func (e *ErrorBuilder) Build() *Error {
	return (*Error)(e)
}

// From creates a new ErrorBuilder from an existing error.
func From(err error) *ErrorBuilder {
	//goland:noinspection GoTypeAssertionOnErrors
	if x, ok := err.(*Error); ok {
		return (*ErrorBuilder)(x.Clone())
	}

	eb := New(err.Error())

	if x, ok := err.(ErrorMessage); ok {
		eb = eb.Msg(x.Message())
	}
	if x, ok := err.(ErrorUserMessage); ok {
		eb = eb.Public(x.UserMessage())
	}
	if x, ok := err.(ErrorHint); ok {
		eb = eb.Hint(x.Hint())
	}
	if x, ok := err.(ErrorCode); ok {
		eb = eb.Code(x.Code())
	}
	if x, ok := err.(ErrorExitCode); ok {
		eb = eb.ExitCode(x.ExitCode())
	}
	if x, ok := err.(ErrorTraceId); ok {
		eb = eb.TraceId(x.TraceId())
	}
	if x, ok := err.(ErrorSpanId); ok {
		eb = eb.SpanId(x.SpanId())
	}
	if x, ok := err.(ErrorTags); ok {
		eb = eb.Tag(x.Tags()...)
	}
	if x, ok := err.(ErrorAttributes); ok {
		eb = eb.AttrsMap(x.Attributes())
	}

	if x, ok := err.(ErrorRelated); ok {
		eb = eb.Related(x.Related()...)
	}

	// Prefer ErrorCauses over other interfaces, since it's more specific, but support all common interfaces for
	// compatibility with other packages.
	if x, ok := err.(ErrorCauses); ok {
		eb = eb.Cause(x.Causes()...)
	} else {
		if x, ok := err.(interface{ Unwrap() []error }); ok {
			eb = eb.Cause(x.Unwrap()...)
		}
		if x, ok := err.(interface{ Unwrap() error }); ok {
			eb = eb.Cause(x.Unwrap())
		}
		if x, ok := err.(interface{ Cause() error }); ok {
			eb = eb.Cause(x.Cause())
		}
	}

	// If the error is a syscall.Errno, use its value as the exit code.
	//goland:noinspection GoTypeAssertionOnErrors
	if x, ok := err.(syscall.Errno); ok {
		eb = eb.ExitCode(int(x))
	}

	return eb
}
