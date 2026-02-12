package ae

// ErrorRecoverable is an interface that should be implemented by error types
// which can indicate whether the error condition is recoverable by returning true
// from ErrorIsRecoverable(), or not recoverable by returning false.
type ErrorRecoverable interface {
	ErrorIsRecoverable() bool
}

// IsRecoverable determines whether the given error is recoverable.
//
// If err is nil, IsRecoverable returns true.
// If err implements the ErrorRecoverable interface, IsRecoverable returns the result of err.ErrorIsRecoverable().
// Otherwise, IsRecoverable defaults to returning true.
func IsRecoverable(err error) bool {
	if err == nil {
		return true
	}

	if ae, ok := err.(ErrorRecoverable); ok {
		return ae.ErrorIsRecoverable()
	}

	return true
}
