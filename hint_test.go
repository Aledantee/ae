package ae_test

import (
	"errors"
	"testing"

	"go.aledante.io/ae"
)

func TestHint_ErrorWithoutInterface(t *testing.T) {
	t.Parallel()

	if got := ae.Hint(errors.New("plain")); got != "" {
		t.Errorf("Hint(plainErr) = %q, want empty string", got)
	}
}

func TestHint_ErrorImplementingInterface(t *testing.T) {
	t.Parallel()

	err := stubErr{msg: "x", hint: "try again later"}
	if got := ae.Hint(err); got != "try again later" {
		t.Errorf("Hint(stubErr) = %q, want %q", got, "try again later")
	}
}

func TestHint_AeBuilderSetsHint(t *testing.T) {
	t.Parallel()

	err := ae.New().Hint("restart the process").Msg("fail")
	if got := ae.Hint(err); got != "restart the process" {
		t.Errorf("Hint on builder = %q, want %q", got, "restart the process")
	}
}
