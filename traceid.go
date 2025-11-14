package ae

// ErrorTraceId defines an interface for errors that can provide a trace ID for distributed tracing.
type ErrorTraceId interface {
	// ErrorTraceId returns the trace ID for distributed tracing.
	// Returns an empty string if no trace ID is set.
	ErrorTraceId() string
}

// TraceId extracts the distributed tracing ID from an error.
// If the error implements ErrorTraceId, returns its TraceId().
// Returns an empty string if err is nil or if the error does not implement ErrorTraceId.
func TraceId(err error) string {
	if err == nil {
		return ""
	}

	if ae, ok := err.(ErrorTraceId); ok {
		return ae.ErrorTraceId()
	}

	return ""
}
