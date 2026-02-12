package ae

// ErrorRecoverable is an interface that should be implemented by error types
// which can indicate whether the error condition is recoverable by returning true
// from ErrorIsRecoverable(), or not recoverable by returning false.
type ErrorRecoverable interface {
	ErrorIsRecoverable() bool
}

// IsRecoverable returns true if the given error is recoverable, and false otherwise.
//
// An error is considered recoverable if:
//   - The error is nil (i.e., no error has occurred).
//   - The error does not implement ErrorRecoverable,
//     or implements it and ErrorIsRecoverable() returns true.
//   - All underlying causes of the error (recursively)
//     are also recoverable.
//
// If any error in the chain implements ErrorRecoverable and its ErrorIsRecoverable() returns false, then the overall error is not recoverable.
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
