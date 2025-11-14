package ae

// ErrorMessage defines an interface for errors that can provide a message.
type ErrorMessage interface {
	// ErrorMessage returns the error message.
	ErrorMessage() string
}

// Message extracts the internal error message from an error.
// If the error implements ErrorMessage, returns its Message().
// Otherwise, returns the error's Error() string.
// Returns an empty string if err is nil.
func Message(err error) string {
	if err == nil {
		return ""
	}

	if ae, ok := err.(ErrorMessage); ok {
		return ae.ErrorMessage()
	}

	return err.Error()
}
