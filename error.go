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
	// parentSpanId links to the parent operation in a trace
	parentSpanId string

	// tags are used for categorizing and filtering errors
	tags map[string]struct{}
	// attrs stores arbitrary key-value pairs for additional context
	attrs map[string]any

	// causes contains the direct causes of this error
	causes []error
	// relatedErrs contains other errors that are related to this error
	relatedErrs []error
	// recoveryErrs contains errors that occurred while handling this error
	recoveryErrs []error
}

// Clone creates a deep copy of the Error.
func (e *Error) Clone() *Error {
	return &Error{
		msg:          e.msg,
		pubMsg:       e.pubMsg,
		hint:         e.hint,
		code:         e.code,
		exitCode:     e.exitCode,
		traceId:      e.traceId,
		spanId:       e.spanId,
		parentSpanId: e.parentSpanId,
		tags:         maps.Clone(e.tags),
		attrs:        maps.Clone(e.attrs),
		causes:       slices.Clone(e.causes),
		relatedErrs:  slices.Clone(e.relatedErrs),
		recoveryErrs: slices.Clone(e.recoveryErrs),
	}
}

// Error implements the error interface for the Error type.
func (e *Error) Error() string {
	return e.msg
}

// Unwrap returns the direct causes of this error.
// It implements the error interface's Unwrap() []error method,
// allowing this error to be used with error wrapping utilities.
func (e *Error) Unwrap() []error {
	return e.causes
}
