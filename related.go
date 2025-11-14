package ae

// ErrorRelated defines an interface for errors that can provide a list of related errors.
// Related errors are those that are not direct causes but are somehow connected to the error,
// including errors that occurred during the handling of the cause(s).
type ErrorRelated interface {
	// ErrorRelated returns a list of errors that are related to the error, but not a direct cause.
	// May also include errors that occurred during the handling of the cause(s).
	// Returns nil if no related errors are set.
	ErrorRelated() []error
}

// Related extracts the list of related errors from an error.
// If the error implements ErrorRelated, returns its Related().
// Returns nil if err is nil or if the error does not implement ErrorRelated.
func Related(err error) []error {
	if err == nil {
		return nil
	}

	if ae, ok := err.(ErrorRelated); ok {
		return ae.ErrorRelated()
	}

	return nil
}
