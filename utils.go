package ae

import (
	"context"
	"fmt"
	"os"
)

// Wrap creates a new error with the given message and wraps the provided error as a cause.
// Returns nil if the provided error is nil.
func Wrap(msg string, err error) error {
	if err == nil {
		return nil
	}

	return From(err).Msg(msg)
}

// ReWrap creates a new error with the provided message and re-wraps all the underlying causes of the given error.
// If the given error is nil or has no causes, it returns nil.
// The resulting error contains the same causes as the input error but with a new top-level message.
func ReWrap(msg string, err error) error {
	causes := Causes(err)
	if causes == nil {
		return nil
	}

	return New().
		Causes(causes).
		Msg(msg)
}

// WrapC creates a new error with the given message and context and wraps the provided error as a cause.
// Returns nil if the provided error is nil.
func WrapC(ctx context.Context, msg string, err error) error {
	if err == nil {
		return nil
	}

	return NewC(ctx).Cause(err).Msg(msg)
}

// Wrapf creates a new error with the given formatted message and wraps the provided errors as causes.
// It filters out any nil errors from the provided list.
// If all provided errors are nil, Wrap returns nil.
// The returned error will have the given message and all non-nil errors as its causes.
func Wrapf(msg string, err error, args ...any) error {
	return Wrap(fmt.Sprintf(msg, args...), err)
}

// WrapCf creates a new error with the given formatted message and context and wraps the provided error as a cause.
// Returns nil if the provided error is nil.
func WrapCf(ctx context.Context, msg string, err error, args ...any) error {
	return NewC(ctx).Cause(err).Msgf(msg, args...)
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

// MsgC creates a new error with the given message and context.
// It is a convenience function that wraps NewC(ctx).Msg(msg).
func MsgC(ctx context.Context, msg string) error {
	return NewC(ctx).Msg(msg)
}

// Msgf creates a new error with the given message.
// It is a convenience function that wraps New().Msgf(msg, args...).
func Msgf(msg string, args ...any) error {
	return New().Msgf(msg, args...)
}

// MsgCf creates a new error with the given message and context.
// It is a convenience function that wraps NewC(ctx).Msgf(msg, args...).
func MsgCf(ctx context.Context, msg string, args ...any) error {
	return NewC(ctx).Msgf(msg, args...)
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
// Does not print anything and exits with code 0 if the error is nil.
func PrintExit(err error, opts ...PrinterOption) {
	Print(err, opts...)
	Exit(err)
}

// Must panics if the provided error is not nil.
// If nil, returns the provided value.
//
// May be used to unwrap errors returned by functions that return (value, error) tuples.
// Example:
//
//	v := Must(SomeFunction())
func Must[T any](v T, err error) T {
	if err != nil {
		panic(fmt.Sprintf("error must not be present: %v", err))
	}

	return v
}

// MustFunc calls the provided function and panics if the returned error is not nil.
// Returns the value returned by the function.
func MustFunc[T any](fn func() (T, error)) T {
	return Must(fn())
}
