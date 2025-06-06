package ae

import (
	"maps"
	"slices"
)

// Error represents a structured error with additional metadata and context.
// It implements the error interface and provides fields for tracing, tagging,
// and error relationships.
type Error struct {
	// msg is the internal error message, typically used for logging
	msg string
	// pubMsg is the public-facing error message, suitable for end users
	pubMsg string
	// hint provides additional guidance or suggestions for resolving the error
	hint string

	// code is an error code that can be used for programmatic error handling
	code string
	// exitCode represents the process exit code associated with this error
	exitCode int

	// traceId is used for distributed tracing
	traceId string
	// spanId identifies the current operation in a trace
	spanId string

	// tags are used for categorizing and filtering errors
	tags map[string]struct{}
	// attrs stores arbitrary key-value pairs for additional context
	attrs map[string]any

	// causes contains the direct causes of this error
	causes []error
	// relatedErrs contains other errors that are related to this error
	relatedErrs []error
}

// Message returns the internal error message, typically used for logging.
func (e *Error) Message() string {
	return e.msg
}

// UserMessage returns the public-facing error message, suitable for end users.
func (e *Error) UserMessage() string {
	return e.pubMsg
}

// Hint returns additional guidance or suggestions for resolving the error.
func (e *Error) Hint() string {
	return e.hint
}

// Code returns the error code that can be used for programmatic error handling.
func (e *Error) Code() string {
	return e.code
}

// ExitCode returns the process exit code associated with this error.
// If the error has no explicit exit code, it returns the highest exit code
// from its causes, or 1 if no causes have an exit code.
func (e *Error) ExitCode() int {
	if e.exitCode > 1 {
		return e.exitCode
	}

	ec := 1
	for _, c := range e.causes {
		if x, ok := c.(ErrorExitCode); ok {
			if x.ExitCode() > ec {
				ec = x.ExitCode()
			}
		}
	}

	return ec
}

// TraceId returns the trace ID used for distributed tracing.
func (e *Error) TraceId() string {
	return e.traceId
}

// SpanId returns the ID that identifies the current operation in a trace.
func (e *Error) SpanId() string {
	return e.spanId
}

// Tags returns a slice of strings containing all tags associated with this error.
func (e *Error) Tags() []string {
	return slices.Collect(maps.Keys(e.tags))
}

// Attributes returns a copy of the arbitrary key-value pairs stored with this error.
func (e *Error) Attributes() map[string]any {
	return maps.Clone(e.attrs)
}

// Causes returns a copy of the direct causes of this error.
func (e *Error) Causes() []error {
	return slices.Clone(e.causes)
}

// Related returns a copy of the errors related to this error.
func (e *Error) Related() []error {
	return slices.Clone(e.relatedErrs)
}

// Clone creates a deep copy of the Error.
// All maps and slices are cloned to ensure the copy is completely independent
// of the original error.
func (e *Error) Clone() *Error {
	return &Error{
		msg:         e.msg,
		pubMsg:      e.pubMsg,
		hint:        e.hint,
		code:        e.code,
		exitCode:    e.exitCode,
		traceId:     e.traceId,
		spanId:      e.spanId,
		tags:        maps.Clone(e.tags),
		attrs:       maps.Clone(e.attrs),
		causes:      slices.Clone(e.causes),
		relatedErrs: slices.Clone(e.relatedErrs),
	}
}

// Error implements the error interface for the Error type.
// It returns the internal error message.
func (e *Error) Error() string {
	return e.msg
}

// Unwrap returns the direct causes of this error.
// It implements the error interface's Unwrap() []error method,
// allowing this error to be used with error wrapping utilities.
func (e *Error) Unwrap() []error {
	return e.causes
}
