package ae

// ErrorCauses defines an interface for errors that can provide a list of underlying causes.
type ErrorCauses interface {
	// ErrorCauses returns a list of errors that caused this error.
	// Returns nil if no causes are set.
	ErrorCauses() []error
}

// Causes extracts the list of underlying causes from an error.
// If the error implements ErrorCauses, returns its Causes().
// If the error implements Unwrap() []error, returns its Unwrap().
// If the error implements Unwrap() error, returns a single-element slice containing its Unwrap().
// Returns nil if err is nil or if the error does not implement any of these interfaces.
func Causes(err error) []error {
	if err == nil {
		return nil
	}

	switch x := err.(type) {
	case ErrorCauses:
		return x.ErrorCauses()
	case interface{ Unwrap() []error }:
		return x.Unwrap()
	case interface{ Unwrap() error }:
		return []error{x.Unwrap()}
	case interface{ Cause() error }:
		return []error{x.Cause()}
	}

	return nil
}
