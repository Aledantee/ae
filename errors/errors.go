package errors

import "github.com/aledantee/ae"

// Join creates a new error that combines multiple errors into a single error.
// It returns nil if no errors are provided.
// The new error will have a message indicating the number of errors that occurred,
// and it will include all the provided errors as causes.
func Join(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	return ae.Newf("%d errors occurred", len(errs)).
		Attrs(errs).
		Build()
}
