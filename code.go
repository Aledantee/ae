package ae

// ErrorCode defines an interface for errors that can provide an error code.
type ErrorCode interface {
	// ErrorCode returns the error code.
	// Returns an empty string if no code is set.
	ErrorCode() string
}

// Code extracts the error code from an error.
// If the error implements ErrorCode, returns its Code().
// Returns an empty string if err is nil or if the error does not implement ErrorCode.
func Code(err error) string {
	if err == nil {
		return ""
	}

	if ae, ok := err.(ErrorCode); ok {
		return ae.ErrorCode()
	}

	return ""
}
