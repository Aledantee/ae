package errors

import (
	stdErrors "errors"
	"strings"

	"go.aledante.io/ae"
)

// New creates a new ae.Ae error with the given message.
func New(msg string) error {
	return ae.New().
		Msg(msg)
}

// Join combines multiple errors into a single error.
// If no errors are provided, it returns nil.
// If only one error is provided, it returns that error directly.
// For multiple errors, it creates a new error with a message containing
// all error messages joined with semicolons and enclosed in square brackets.
func Join(errs ...error) error {
	var filtered []error
	for _, err := range errs {
		if err != nil {
			filtered = append(filtered, err)
		}
	}

	switch len(errs) {
	case 0:
		return nil
	case 1:
		return filtered[0]
	default:
		var sb strings.Builder
		sb.WriteRune('[')

		for i, err := range filtered {
			if i > 0 {
				sb.WriteString("; ")
			}

			sb.WriteString(err.Error())
		}

		sb.WriteRune(']')
		return ae.New().
			Causes(filtered).
			Msg(sb.String())
	}
}

// Is reports whether any error in err's chain matches target.
// It is a proxy to the standard errors.Is.
func Is(err, target error) bool {
	return stdErrors.Is(err, target)
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true. Otherwise, it returns false.
// It is a proxy to the standard errors.As.
func As(err error, target any) bool {
	//goland:noinspection GoErrorsAs
	return stdErrors.As(err, target)
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error. Otherwise, Unwrap returns nil.
// It is a proxy to the standard errors.Unwrap.
func Unwrap(err error) error {
	return stdErrors.Unwrap(err)
}
