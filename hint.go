package ae

// ErrorHint defines an interface for errors that can provide a hint for resolution.
type ErrorHint interface {
	// ErrorHint returns a hint for resolving the error.
	// Returns an empty string if no hint is set.
	ErrorHint() string
}

func Hint(err error) string {
	if ae, ok := err.(ErrorHint); ok {
		return ae.ErrorHint()
	}

	return ""
}
