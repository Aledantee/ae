package ae_test

import (
	"errors"
	"testing"

	"go.aledante.io/ae"
)

func TestCauses_NilError(t *testing.T) {
	t.Parallel()

	if got := ae.Causes(nil); got != nil {
		t.Errorf("Causes(nil) = %v, want nil", got)
	}
}

func TestCauses_ErrorWithoutRelevantInterface(t *testing.T) {
	t.Parallel()

	// Bare stdlib error: no ErrorCauses, Unwrap, or Cause implementation.
	if got := ae.Causes(errors.New("plain")); got != nil {
		t.Errorf("Causes(plainErr) = %v, want nil", got)
	}
}

func TestCauses_ErrorCausesInterface(t *testing.T) {
	t.Parallel()

	c1 := errors.New("c1")
	c2 := errors.New("c2")
	err := stubErr{msg: "x", causes: []error{c1, c2}}
	got := ae.Causes(err)
	if len(got) != 2 || got[0] != c1 || got[1] != c2 {
		t.Errorf("Causes(ErrorCauses) = %v, want [c1 c2]", got)
	}
}

func TestCauses_UnwrapMultiInterface(t *testing.T) {
	t.Parallel()

	c1 := errors.New("c1")
	c2 := errors.New("c2")
	err := multiUnwrapErr{msg: "multi", errs: []error{c1, c2}}
	got := ae.Causes(err)
	if len(got) != 2 || got[0] != c1 || got[1] != c2 {
		t.Errorf("Causes(Unwrap []error) = %v, want [c1 c2]", got)
	}
}

func TestCauses_UnwrapSingleInterface(t *testing.T) {
	t.Parallel()

	inner := errors.New("inner")
	err := singleUnwrapErr{msg: "outer", inner: inner}
	got := ae.Causes(err)
	if len(got) != 1 || got[0] != inner {
		t.Errorf("Causes(Unwrap error) = %v, want single-element slice containing inner", got)
	}
}

func TestCauses_PkgErrorsStyleCauseInterface(t *testing.T) {
	t.Parallel()

	inner := errors.New("root")
	err := pkgErrorsStyleErr{msg: "wrapper", inner: inner}
	got := ae.Causes(err)
	if len(got) != 1 || got[0] != inner {
		t.Errorf("Causes(Cause error) = %v, want single-element slice containing inner", got)
	}
}

func TestCauses_ErrorCausesWinsOverUnwrap(t *testing.T) {
	t.Parallel()

	// stubErr implements ErrorCauses. The switch order in Causes(err) documents
	// that ErrorCauses is checked first.
	c := errors.New("expected")
	err := stubErr{msg: "x", causes: []error{c}}
	got := ae.Causes(err)
	if len(got) != 1 || got[0] != c {
		t.Errorf("Causes precedence: got %v, want [%v]", got, c)
	}
}
