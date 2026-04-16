package ae_test

import (
	"errors"
	"testing"

	"go.aledante.io/ae"
)

func TestIsRecoverable_NilIsRecoverable(t *testing.T) {
	t.Parallel()

	if !ae.IsRecoverable(nil) {
		t.Error("IsRecoverable(nil) = false, want true (nil is treated as no error)")
	}
}

func TestIsRecoverable_PlainErrorIsRecoverable(t *testing.T) {
	t.Parallel()

	// A plain error that does not implement ErrorRecoverable is recoverable
	// by default, per the documented contract.
	if !ae.IsRecoverable(errors.New("plain")) {
		t.Error("IsRecoverable(plainErr) = false, want true")
	}
}

func TestIsRecoverable_BuilderDefaultsToRecoverable(t *testing.T) {
	t.Parallel()

	// New() initialises recoverable=true so freshly-built errors are
	// considered recoverable without an explicit .Recoverable(true) call.
	err := ae.New().Msg("x")
	if !ae.IsRecoverable(err) {
		t.Error("fresh Ae reported as not recoverable")
	}
}

func TestIsRecoverable_FatalBuilderIsNotRecoverable(t *testing.T) {
	t.Parallel()

	err := ae.New().Fatal().Msg("x")
	if ae.IsRecoverable(err) {
		t.Error("Fatal() Ae reported as recoverable")
	}
}

func TestIsRecoverable_RecoverableFalseIsNotRecoverable(t *testing.T) {
	t.Parallel()

	err := ae.New().Recoverable(false).Msg("x")
	if ae.IsRecoverable(err) {
		t.Error("Recoverable(false) Ae reported as recoverable")
	}
}

func TestIsRecoverable_UnrecoverableCauseMakesChainUnrecoverable(t *testing.T) {
	t.Parallel()

	// Outer is recoverable by default but wraps a cause marked fatal.
	// Per the docstring, any non-recoverable error anywhere in the chain
	// makes the whole chain non-recoverable.
	inner := ae.New().Fatal().Msg("fatal")
	outer := ae.New().Cause(inner).Msg("outer")

	if ae.IsRecoverable(outer) {
		t.Error("chain with fatal cause reported as recoverable")
	}
}

func TestIsRecoverable_AllCausesRecoverable(t *testing.T) {
	t.Parallel()

	inner1 := errors.New("plain")
	inner2 := ae.New().Msg("default-recoverable")
	outer := ae.New().Cause(inner1, inner2).Msg("outer")

	if !ae.IsRecoverable(outer) {
		t.Error("chain of recoverable errors reported as not recoverable")
	}
}

// stubRecoverableErr lets tests exercise the ErrorRecoverable interface on a
// non-ae error type without needing to build a full stubErr.
type stubRecoverableErr struct {
	msg         string
	recoverable bool
}

func (s stubRecoverableErr) Error() string            { return s.msg }
func (s stubRecoverableErr) ErrorIsRecoverable() bool { return s.recoverable }

func TestIsRecoverable_ExternalErrorTypeHonoured(t *testing.T) {
	t.Parallel()

	if ae.IsRecoverable(stubRecoverableErr{msg: "x", recoverable: false}) {
		t.Error("external error with ErrorIsRecoverable()=false reported as recoverable")
	}
	if !ae.IsRecoverable(stubRecoverableErr{msg: "x", recoverable: true}) {
		t.Error("external error with ErrorIsRecoverable()=true reported as not recoverable")
	}
}
