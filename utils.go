package ae

// Wrap creates a new error with the given message and wraps the provided errors as causes.
// It filters out any nil errors from the provided list.
// If all provided errors are nil, Wrap returns nil.
// The returned error will have the given message and all non-nil errors as its causes.
func Wrap(msg string, errs ...error) error {
	var filtered []error
	for _, err := range errs {
		if err != nil {
			filtered = append(filtered, err)
		}
	}

	if len(filtered) == 0 {
		return nil
	}

	return New().
		Causes(filtered).
		Msg(msg)
}

// Msg creates a new error with the given message.
// It is a convenience function that wraps New().Msg(msg).
func Msg(msg string) error {
	return New().Msg(msg)
}
