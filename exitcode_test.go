package ae_test

import (
	"errors"
	"testing"

	"go.aledante.io/ae"
)

func TestExitCode_NilError(t *testing.T) {
	t.Parallel()

	if got := ae.ExitCode(nil); got != 0 {
		t.Errorf("ExitCode(nil) = %d, want 0", got)
	}
}

func TestExitCode_PlainErrorDefaultsToOne(t *testing.T) {
	t.Parallel()

	if got := ae.ExitCode(errors.New("plain")); got != 1 {
		t.Errorf("ExitCode(plainErr) = %d, want 1 (documented default for non-nil error)", got)
	}
}

func TestExitCode_UsesOwnExitCodeWhenPositive(t *testing.T) {
	t.Parallel()

	err := stubErr{msg: "x", exitCode: 42}
	if got := ae.ExitCode(err); got != 42 {
		t.Errorf("ExitCode(stubErr) = %d, want 42", got)
	}
}

func TestExitCode_RecursivelyMaxOverCauses(t *testing.T) {
	t.Parallel()

	// No exit code on the outer error. Inner causes have 3 and 7; expect 7.
	inner1 := stubErr{msg: "i1", exitCode: 3}
	inner2 := stubErr{msg: "i2", exitCode: 7}
	outer := ae.New().Cause(inner1, inner2).Msg("outer")

	if got := ae.ExitCode(outer); got != 7 {
		t.Errorf("ExitCode recursive max = %d, want 7", got)
	}
}

func TestExitCode_DeepRecursion(t *testing.T) {
	t.Parallel()

	leaf := stubErr{msg: "leaf", exitCode: 9}
	mid := ae.New().Cause(leaf).Msg("mid")
	top := ae.New().Cause(mid).Msg("top")

	if got := ae.ExitCode(top); got != 9 {
		t.Errorf("ExitCode deep recursive = %d, want 9", got)
	}
}

func TestBuilder_ExitCodeStoresPositiveOnly(t *testing.T) {
	t.Parallel()

	// Positive is stored.
	err := ae.New().ExitCode(5).Msg("x")
	if got := ae.ExitCode(err); got != 5 {
		t.Errorf("ExitCode after ExitCode(5) = %d, want 5", got)
	}

	// Zero and negative are ignored per the docstring "Only positive values are stored."
	err = ae.New().ExitCode(5).ExitCode(0).ExitCode(-1).Msg("x")
	if got := ae.ExitCode(err); got != 5 {
		t.Errorf("ExitCode after ExitCode(0/-1) overwrote = %d, want 5", got)
	}
}

// TestAe_ErrorExitCodeInterfaceRecurses exercises the ErrorExitCode interface
// contract: "If the error does not have an associated exit code, the highest
// exit code of all recursive causes is returned." *Ae currently returns only
// its own exitCode field; expected to fail until the recursive walk is added.
func TestAe_ErrorExitCodeInterfaceRecurses(t *testing.T) {
	t.Parallel()

	leaf := stubErr{msg: "leaf", exitCode: 11}
	outer := ae.New().Cause(leaf).Msg("outer")

	got, ok := outer.(ae.ErrorExitCode)
	if !ok {
		t.Fatalf("*ae.Ae does not implement ErrorExitCode")
	}
	if ec := got.ErrorExitCode(); ec != 11 {
		t.Errorf("Ae.ErrorExitCode() = %d, want 11", ec)
	}
}
