package ae

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
