package ae

import (
	"maps"
	"slices"
	"strings"
	"time"
)

// Ae represents an error type that implements the standard error interface along with
// additional interfaces defined in the ae package for enhanced error handling.
type Ae struct {
	// msg is the internal error message, typically used for logging and debugging
	msg string
	// userMsg is a user-friendly error message that can be safely displayed to end users
	userMsg string
	// hint provides additional guidance or suggestions for resolving the error
	hint string

	// timestamp is the time the error occurred
	timestamp time.Time

	// code is an error code that can be used for programmatic error handling
	code string
	// exitCode represents the process exit code that should be used when this error occurs
	exitCode int

	// traceId is used for distributed tracing to correlate related operations
	traceId string
	// spanId identifies a specific operation within a trace
	spanId string

	// tags are used to categorize and filter errors
	tags map[string]struct{}
	// attributes provide additional context-specific information about the error
	attributes map[string]any

	// joined indicates whether this error represents a collection of joined errors
	joined bool
	// causes contains the underlying errors that led to this error
	causes []error
	// related contains errors that are related to this error, but not a direct cause
	// also includes errors that occurred during the handling of the cause(s)
	related []error

	// stacks contains the stack traces associated with this error
	stacks []*Stack
}

// New creates a new error with the given message.
// The error will have a timestamp set to the current time and will implement
// all the error interfaces defined in the ae package.
func New(msg string) error {
	return newInternal(msg)
}

// Message returns the internal error message.
func (a *Ae) Message() string {
	return a.msg
}

// UserMessage returns the user-friendly error message.
func (a *Ae) UserMessage() string {
	return a.userMsg
}

// Hint returns additional guidance for resolving the error.
func (a *Ae) Hint() string {
	return a.hint
}

// Timestamp returns the timestamp of the error.
func (a *Ae) Timestamp() time.Time {
	return a.timestamp
}

// Code returns the error code.
func (a *Ae) Code() string {
	return a.code
}

// ExitCode returns the process exit code associated with this error.
func (a *Ae) ExitCode() int {
	return a.exitCode
}

// TraceId returns the distributed tracing ID.
func (a *Ae) TraceId() string {
	return a.traceId
}

// SpanId returns the operation span ID.
func (a *Ae) SpanId() string {
	return a.spanId
}

// Tags returns a slice of all tags associated with this error.
func (a *Ae) Tags() []string {
	return slices.Collect(maps.Keys(a.tags))
}

// Attributes returns a copy of the error's attributes map.
func (a *Ae) Attributes() map[string]any {
	return maps.Clone(a.attributes)
}

// IsJoined returns true if this error represents a collection of joined errors.
func (a *Ae) IsJoined() bool {
	return a.joined
}

// Causes returns a copy of the underlying errors that caused this error.
func (a *Ae) Causes() []error {
	return slices.Clone(a.causes)
}

// Related returns a copy of the errors that are related to this error, but not a direct cause.
// Also includes errors that occurred during the handling of the cause(s).
func (a *Ae) Related() []error {
	return slices.Clone(a.related)
}

// Stacks returns a copy of the stack traces associated with this error.
func (a *Ae) Stacks() []*Stack {
	return slices.Clone(a.stacks)
}

// Error implements the error interface by returning a string representation of the error.
// It includes the main error message and any underlying causes.
func (a *Ae) Error() string {
	var errMsg strings.Builder
	errMsg.WriteString(a.msg)

	if len(a.causes) > 0 {
		errMsg.WriteString(": ")

		if len(a.causes) == 1 {
			errMsg.WriteString(a.causes[0].Error())
		} else {
			errMsg.WriteString("[")
			for i, cause := range a.causes {
				if i > 0 {
					errMsg.WriteString("; ")
				}
				errMsg.WriteString(cause.Error())
			}
			errMsg.WriteString("]")
		}
	}

	return errMsg.String()
}

// Unwrap returns the underlying errors that caused this error.
// This implements the errors.Unwrap interface.
func (a *Ae) Unwrap() []error {
	return a.Causes()
}

// Print writes the formatted error to standard output using the provided printer options.
func (a *Ae) Print(opts ...PrinterOption) {
	NewPrinter(opts...).Print(a)
}

// Prints returns a string representation of the error using the provided printer options.
func (a *Ae) Prints(opts ...PrinterOption) string {
	return NewPrinter(opts...).Prints(a)
}

// newInternal creates a new Ae error with the given message.
// It initializes the error with the current timestamp and empty maps for tags and attributes.
func newInternal(msg string) *Ae {
	return &Ae{
		msg:        msg,
		timestamp:  time.Now(),
		tags:       make(map[string]struct{}),
		attributes: make(map[string]any),
	}
}
