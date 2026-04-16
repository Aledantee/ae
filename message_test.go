package ae_test

import (
	"errors"
	"testing"

	"go.aledante.io/ae"
)

func TestMessage_NilError(t *testing.T) {
	t.Parallel()

	if got := ae.Message(nil); got != "" {
		t.Errorf("Message(nil) = %q, want empty string", got)
	}
}

func TestMessage_ErrorImplementingInterface(t *testing.T) {
	t.Parallel()

	// The internal message and the Error() string are distinct on purpose: ae
	// documents that Message should prefer the ErrorMessage interface.
	err := stubErr{msg: "internal"}
	if got := ae.Message(err); got != "internal" {
		t.Errorf("Message(stubErr) = %q, want %q", got, "internal")
	}
}

func TestMessage_FallsBackToErrorString(t *testing.T) {
	t.Parallel()

	err := errors.New("boom")
	if got := ae.Message(err); got != "boom" {
		t.Errorf("Message(stdlibErr) = %q, want %q (fallback to err.Error())", got, "boom")
	}
}

func TestMessage_AeBuilderReturnsInternalMessage(t *testing.T) {
	t.Parallel()

	err := ae.New().Msg("hello")
	if got := ae.Message(err); got != "hello" {
		t.Errorf("Message on builder = %q, want %q", got, "hello")
	}
}
