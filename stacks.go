package ae

// ErrorStacks defines an interface for errors that can provide a list of stack traces.
type ErrorStacks interface {
	// ErrorStacks returns a list of stack traces associated with the error, one for each goroutine.
	// Returns nil if no stack traces are set.
	ErrorStacks() []*Stack
}

// Stacks extracts the list of stack traces from an error.
// If the error implements ErrorStacks, returns its Stacks().
// Returns nil if err is nil or if the error does not implement ErrorStacks.
func Stacks(err error) []*Stack {
	if err == nil {
		return nil
	}

	if ae, ok := err.(ErrorStacks); ok {
		return ae.ErrorStacks()
	}

	return nil
}
