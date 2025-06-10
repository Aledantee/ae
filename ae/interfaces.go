package ae

import "time"

// ErrorMessage defines an interface for errors that can provide a message.
type ErrorMessage interface {
	// Message returns the error message.
	Message() string
}

// ErrorUserMessage defines an interface for errors that can an error message for end-users.
type ErrorUserMessage interface {
	// UserMessage returns an error message for end-users.
	// Returns an empty string if no end-user message is set.
	UserMessage() string
}

// ErrorHint defines an interface for errors that can provide a hint for resolution.
type ErrorHint interface {
	// Hint returns a hint for resolving the error.
	// Returns an empty string if no hint is set.
	Hint() string
}

// ErrorTimestamp defines an interface for errors that can provide a timestamp.
type ErrorTimestamp interface {
	// Timestamp returns the timestamp of the error.
	// Returns zero time if no timestamp is set.
	Timestamp() time.Time
}

// ErrorCode defines an interface for errors that can provide an error code.
type ErrorCode interface {
	// Code returns the error code.
	// Returns an empty string if no code is set.
	Code() string
}

// ErrorExitCode defines an interface for errors that can provide an exit code.
type ErrorExitCode interface {
	// ExitCode returns the exit code associated with the error.
	// If the error does not have an associated exit code, the highest exit code of all recursive causes is returned.
	ExitCode() int
}

// ErrorTraceId defines an interface for errors that can provide a trace ID for distributed tracing.
type ErrorTraceId interface {
	// TraceId returns the trace ID for distributed tracing.
	// Returns an empty string if no trace ID is set.
	TraceId() string
}

// ErrorSpanId defines an interface for errors that can provide a span ID for distributed tracing.
type ErrorSpanId interface {
	// SpanId returns the span ID for distributed tracing.
	// Returns an empty string if no span ID is set.
	SpanId() string
}

// ErrorTags defines an interface for errors that can provide a list of tags.
type ErrorTags interface {
	// Tags returns a list of tags associated with the error.
	// Returns nil if no tags are set.
	Tags() []string
}

// ErrorAttributes defines an interface for errors that can provide a map of attributes.
type ErrorAttributes interface {
	// Attributes returns a map of attributes associated with the error.
	// Returns an empty map non-nil if no attributes are set.
	Attributes() map[string]any
}

// ErrorCauses defines an interface for errors that can provide a list of underlying causes.
type ErrorCauses interface {
	// Causes returns a list of errors that caused this error.
	// Returns nil if no causes are set.
	Causes() []error
}

// ErrorJoined defines an interface for errors that represent a collection of independent errors joined together.
// Unlike wrapped errors where one error contains another, joined errors treat all causes as peers.
// A joined error typically does not have its own message, hint, code, or exit code - it serves as a container
// for multiple independent errors that occurred in parallel or are otherwise related but not causally linked.
type ErrorJoined interface {
	ErrorCauses

	// IsJoined returns true if this error represents a collection of joined errors.
	// Returns false if this is a single error or if the error was not created through joining.
	// This can be used to distinguish between errors that wrap causes and those that join them.
	IsJoined() bool
}

// ErrorRelated defines an interface for errors that can provide a list of related errors.
// Related errors are those that are not direct causes but are somehow connected to the error,
// including errors that occurred during the handling of the cause(s).
type ErrorRelated interface {
	// Related returns a list of errors that are related to the error, but not a direct cause.
	// May also include errors that occurred during the handling of the cause(s).
	// Returns nil if no related errors are set.
	Related() []error
}

// ErrorStacks defines an interface for errors that can provide a list of stack traces.
type ErrorStacks interface {
	// Stacks returns a list of stack traces associated with the error, one for each goroutine.
	// Returns nil if no stack traces are set.
	Stacks() []*Stack
}
