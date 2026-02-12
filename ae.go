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
	// recoverable indicates whether the error is recoverable
	recoverable bool

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

	// causes contains the underlying errors that led to this error
	causes []error
	// related contains errors that are related to this error, but not a direct cause
	// also includes errors that occurred during the handling of the cause(s)
	related []error

	// stacks contains the stack traces associated with this error
	stacks []*Stack
}

// ErrorMessage returns the internal error message.
func (a Ae) ErrorMessage() string {
	return a.msg
}

// ErrorUserMessage returns the user-friendly error message.
func (a Ae) ErrorUserMessage() string {
	return a.userMsg
}

// ErrorHint returns additional guidance for resolving the error.
func (a Ae) ErrorHint() string {
	return a.hint
}

// ErrorIsRecoverable returns whether the error is recoverable.
func (a Ae) ErrorIsRecoverable() bool {
	return a.recoverable
}

// ErrorTimestamp returns the timestamp of the error.
func (a Ae) ErrorTimestamp() time.Time {
	return a.timestamp
}

// ErrorCode returns the error code.
func (a Ae) ErrorCode() string {
	return a.code
}

// ErrorExitCode returns the process exit code associated with this error.
func (a Ae) ErrorExitCode() int {
	return a.exitCode
}

// ErrorTraceId returns the distributed tracing ID.
func (a Ae) ErrorTraceId() string {
	return a.traceId
}

// ErrorSpanId returns the operation span ID.
func (a Ae) ErrorSpanId() string {
	return a.spanId
}

// ErrorTags returns a slice of all tags associated with this error.
func (a Ae) ErrorTags() []string {
	return slices.Collect(maps.Keys(a.tags))
}

// ErrorAttributes returns a copy of the error's attributes map.
func (a Ae) ErrorAttributes() map[string]any {
	return maps.Clone(a.attributes)
}

// ErrorCauses returns a copy of the underlying errors that caused this error.
func (a Ae) ErrorCauses() []error {
	return slices.Clone(a.causes)
}

// ErrorRelated returns a copy of the errors that are related to this error, but not a direct cause.
// Also includes errors that occurred during the handling of the cause(s).
func (a Ae) ErrorRelated() []error {
	return slices.Clone(a.related)
}

// ErrorStacks returns a copy of the stack traces associated with this error.
func (a Ae) ErrorStacks() []*Stack {
	return slices.Clone(a.stacks)
}

// Error implements the error interface by returning a string representation of the error.
// It includes the main error message and any underlying causes.
func (a Ae) Error() string {
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
func (a Ae) Unwrap() []error {
	return a.ErrorCauses()
}

// Print writes the formatted error to standard output using the provided printer options.
func (a Ae) Print(opts ...PrinterOption) {
	NewPrinter(opts...).Print(a)
}

// Prints returns a string representation of the error using the provided printer options.
func (a Ae) Prints(opts ...PrinterOption) string {
	return NewPrinter(opts...).Prints(a)
}

// clone creates and returns a deep copy of the Ae instance and its associated fields.
func (a Ae) clone() Ae {
	cpy := a

	cpy.tags = maps.Clone(a.tags)
	cpy.attributes = maps.Clone(a.attributes)
	cpy.causes = slices.Clone(a.causes)
	cpy.related = slices.Clone(a.related)
	cpy.stacks = slices.Clone(a.stacks)

	return cpy
}
