package ae

// ErrorRecoverable is an interface that should be implemented by error types
// which can indicate whether the error condition is recoverable by returning true
// from ErrorIsRecoverable(), or not recoverable by returning false.
type ErrorRecoverable interface {
	ErrorIsRecoverable() bool
}

func IsRecoverable(err error) bool {
	if err == nil {
		return true
	}

	if ae, ok := err.(ErrorRecoverable); ok && !ae.ErrorIsRecoverable() {
		return false
	}

	for _, cause := range Causes(err) {
		if !IsRecoverable(cause) {
			return false
		}
	}

	return true
}
