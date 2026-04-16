package ae_test

import (
	"errors"
	"testing"

	"go.aledante.io/ae"
)

func TestCode_NilError(t *testing.T) {
	t.Parallel()

	if got := ae.Code(nil); got != "" {
		t.Errorf("Code(nil) = %q, want empty string", got)
	}
}

func TestCode_ErrorWithoutInterface(t *testing.T) {
	t.Parallel()

	if got := ae.Code(errors.New("plain")); got != "" {
		t.Errorf("Code(plainErr) = %q, want empty string (not implementing ErrorCode)", got)
	}
}

func TestCode_ErrorImplementingInterface(t *testing.T) {
	t.Parallel()

	err := stubErr{msg: "x", code: "DB_DOWN"}
	if got := ae.Code(err); got != "DB_DOWN" {
		t.Errorf("Code(stubErr) = %q, want %q", got, "DB_DOWN")
	}
}

func TestCode_AeBuilderSetsCode(t *testing.T) {
	t.Parallel()

	err := ae.New().Code("AUTH_FAIL").Msg("oops")
	if got := ae.Code(err); got != "AUTH_FAIL" {
		t.Errorf("Code on ae builder = %q, want %q", got, "AUTH_FAIL")
	}
}
