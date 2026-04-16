package errors_test

import (
	stdErrors "errors"
	"strings"
	"testing"

	aeerrors "go.aledante.io/ae/errors"
)

func TestNew_ReturnsErrorWithMessage(t *testing.T) {
	t.Parallel()

	err := aeerrors.New("boom")
	if err == nil {
		t.Fatal("New returned nil")
	}
	if err.Error() != "boom" {
		t.Errorf("Error() = %q, want %q", err.Error(), "boom")
	}
}

// TestJoin_NoArgsReturnsNil asserts the documented behavior: "If no errors are
// provided, it returns nil." The current implementation's switch on the
// unfiltered len(errs) correctly handles this case.
func TestJoin_NoArgsReturnsNil(t *testing.T) {
	t.Parallel()

	if got := aeerrors.Join(); got != nil {
		t.Errorf("Join() = %v, want nil", got)
	}
}

// TestJoin_AllNilReturnsNil asserts that Join treats "no errors provided" as
// "no non-nil errors". The current implementation does not filter before the
// switch, so Join(nil) hits case 1 and panics indexing into an empty filtered
// slice; Join(nil, nil) builds a bracketed empty message. Expected to fail
// until the filtering is moved above the switch.
func TestJoin_AllNilReturnsNil(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Join(nil) panicked: %v", r)
		}
	}()

	if got := aeerrors.Join(nil, nil); got != nil {
		t.Errorf("Join(nil, nil) = %v, want nil", got)
	}
}

func TestJoin_SingleErrorReturnedDirectly(t *testing.T) {
	t.Parallel()

	e := stdErrors.New("only")
	got := aeerrors.Join(e)
	if got != e {
		t.Errorf("Join(e) = %v, want the original error returned directly", got)
	}
}

func TestJoin_MultipleErrorsBracketedWithCauses(t *testing.T) {
	t.Parallel()

	e1 := stdErrors.New("one")
	e2 := stdErrors.New("two")
	e3 := stdErrors.New("three")

	got := aeerrors.Join(e1, e2, e3)
	if got == nil {
		t.Fatal("Join returned nil for non-empty input")
	}

	msg := got.Error()
	if !strings.Contains(msg, "one") || !strings.Contains(msg, "two") || !strings.Contains(msg, "three") {
		t.Errorf("Join error = %q, want to contain all sub-messages", msg)
	}
	// Doc says the joined message is bracketed and semicolon-separated.
	if !strings.Contains(msg, "[") || !strings.Contains(msg, "]") {
		t.Errorf("Join error = %q, want bracketed form", msg)
	}

	// And errors.Is walks the causes via the underlying Ae's Unwrap() []error.
	if !stdErrors.Is(got, e1) {
		t.Errorf("errors.Is(joined, e1) = false, want true")
	}
}

func TestIs_ProxiesToStdErrors(t *testing.T) {
	t.Parallel()

	target := stdErrors.New("target")
	wrapped := aeerrors.Join(target, stdErrors.New("other"))

	if !aeerrors.Is(wrapped, target) {
		t.Errorf("Is(wrapped, target) = false, want true")
	}
	if aeerrors.Is(wrapped, stdErrors.New("never")) {
		t.Errorf("Is with unrelated target returned true")
	}
}

func TestAs_ProxiesToStdErrors(t *testing.T) {
	t.Parallel()

	sentinel := &wrappableErr{msg: "sentinel"}
	wrapped := aeerrors.Join(sentinel, stdErrors.New("noise"))

	var target *wrappableErr
	if !aeerrors.As(wrapped, &target) {
		t.Errorf("As(wrapped, &target) = false, want true")
	}
	if target == nil || target.msg != "sentinel" {
		t.Errorf("As populated target = %+v, want sentinel", target)
	}
}

type wrappableErr struct{ msg string }

func (w *wrappableErr) Error() string { return w.msg }

func TestUnwrap_ProxiesToStdErrors(t *testing.T) {
	t.Parallel()

	// A manually wrapped stdlib error so Unwrap has something to return.
	inner := stdErrors.New("inner")
	outer := newWrapper("outer", inner)

	if got := aeerrors.Unwrap(outer); got != inner {
		t.Errorf("Unwrap(outer) = %v, want inner", got)
	}
}

type wrapper struct {
	msg   string
	inner error
}

func (w *wrapper) Error() string { return w.msg + ": " + w.inner.Error() }
func (w *wrapper) Unwrap() error { return w.inner }

func newWrapper(msg string, inner error) *wrapper { return &wrapper{msg: msg, inner: inner} }
