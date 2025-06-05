package ae

import (
	"fmt"
	"syscall"
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

// Tag adds a tag to categorize or filter the error.
func (e *ErrorBuilder) Tag(tag string) *ErrorBuilder {
	e.tags[tag] = struct{}{}
	return e
}

// Attr sets a single attribute with the given key and value.
func (e *ErrorBuilder) Attr(key string, value any) *ErrorBuilder {
	e.attrs[key] = value
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

		e.attrs[key] = kv[i+1]
	}

	return e
}

// Cause adds one or more errors as direct causes of this error.
func (e *ErrorBuilder) Cause(errs ...error) *ErrorBuilder {
	e.causes = append(e.causes, errs...)
	return e
}

// Related adds one or more errors that are related to this error.
func (e *ErrorBuilder) Related(errs ...error) *ErrorBuilder {
	e.relatedErrs = append(e.relatedErrs, errs...)
	return e
}

// Recovery adds one or more errors that occurred while handling this error.
func (e *ErrorBuilder) Recovery(errs ...error) *ErrorBuilder {
	e.recoveryErrs = append(e.recoveryErrs, errs...)
	return e
}

// Build converts the ErrorBuilder to an Error.
func (e *ErrorBuilder) Build() *Error {
	return (*Error)(e)
}

// From creates a new ErrorBuilder from an existing error.
// It handles various error types and interfaces:
// - *Error: Creates a builder based on the error
// - Unwrap() []error: Adds unwrapped errors as causes
// - Unwrap() error: Adds the unwrapped error as a cause
// - Cause() error: Adds the cause as a cause
// - syscall.Errno: Sets the exit code to the error number
func From(err error) *ErrorBuilder {
	if x, ok := err.(*Error); ok {
		return (*ErrorBuilder)(x.Clone())
	}

	eb := New(err.Error())

	if x, ok := err.(interface{ Unwrap() []error }); ok {
		eb = eb.Cause(x.Unwrap()...)
	}
	if x, ok := err.(interface{ Unwrap() error }); ok {
		eb = eb.Cause(x.Unwrap())
	}
	if x, ok := err.(interface{ Cause() error }); ok {
		eb = eb.Cause(x.Cause())
	}
	if x, ok := err.(syscall.Errno); ok {
		eb = eb.ExitCode(int(x))
	}

	return eb
}
