package ae

import (
	"fmt"
	"os"
)

// Wrap creates a new error with the given message and wraps the provided error as a cause.
// Returns nil if the provided error is nil.
func Wrap(msg string, err error) error {
	if err == nil {
		return nil
	}

	return New().
		Cause(err).
		Msg(msg)
}

// Wrapf creates a new error with the given formatted message and wraps the provided errors as causes.
// It filters out any nil errors from the provided list.
// If all provided errors are nil, Wrap returns nil.
// The returned error will have the given message and all non-nil errors as its causes.
func Wrapf(msg string, err error, args ...any) error {
	return Wrap(fmt.Sprintf(msg, args...), err)
}

// WrapMany creates a new error with the given message and wraps the provided errors as causes.
// It filters out any nil errors from the provided list.
// If all provided errors are nil, Wrap returns nil.
// The returned error will have the given message and all non-nil errors as its causes.
func WrapMany(msg string, errs ...error) error {
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

// Msgf creates a new error with the given message.
// It is a convenience function that wraps New().Msg(msg).
func Msgf(msg string, args ...any) error {
	return Msg(fmt.Sprintf(msg, args...))
}

// Exit exits the program with the exit code returned by ExitCode.
// Does nothing if the error is nil.
func Exit(err error) {
	if err == nil {
		return
	}

	os.Exit(ExitCode(err))
}

// PrintExit prints the error to stderr and exits the program with the exit code returned by ExitCode.
// Does nothing if the error is nil.
func PrintExit(err error) {
	Print(err)
	Exit(err)
}
