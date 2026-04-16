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
// Nil entries are filtered before the combination is decided:
//   - If all inputs are nil (or the list is empty), returns nil.
//   - If exactly one non-nil error is supplied, returns it directly.
//   - Otherwise, creates an ae error whose message joins every sub-message
//     with semicolons inside square brackets and whose causes are the
//     surviving non-nil errors.
func Join(errs ...error) error {
	var filtered []error
	for _, err := range errs {
		if err != nil {
			filtered = append(filtered, err)
		}
	}

	switch len(filtered) {
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
