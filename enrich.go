package ae

// AddRelated creates a new error that includes the original error and adds related errors to it.
// It preserves all properties of the original error and appends the provided related errors.
// The error message from the original error is used for the new error.
func AddRelated(err error, related ...error) error {
	return From(err).
		Related(related...).
		Msg(Message(err))
}
