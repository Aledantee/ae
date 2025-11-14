package ae

// ErrorUserMessage defines an interface for errors that can an error message for end-users.
type ErrorUserMessage interface {
	// ErrorUserMessage returns an error message for end-users.
	// Returns an empty string if no end-user message is set.
	ErrorUserMessage() string
}

// UserMessage extracts the user-friendly error message from an error.
// If the error implements ErrorUserMessage, returns its UserMessage().
// Returns an empty string if err is nil or if the error does not implement ErrorUserMessage.
func UserMessage(err error) string {
	if err == nil {
		return ""
	}

	if ae, ok := err.(ErrorUserMessage); ok {
		return ae.ErrorUserMessage()
	}

	return ""
}
