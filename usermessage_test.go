package ae_test

import (
	"errors"
	"testing"

	"go.aledante.io/ae"
)

func TestUserMessage_NilError(t *testing.T) {
	t.Parallel()

	if got := ae.UserMessage(nil); got != "" {
		t.Errorf("UserMessage(nil) = %q, want empty string", got)
	}
}

func TestUserMessage_ErrorWithoutInterface(t *testing.T) {
	t.Parallel()

	if got := ae.UserMessage(errors.New("plain")); got != "" {
		t.Errorf("UserMessage(plainErr) = %q, want empty string", got)
	}
}

func TestUserMessage_ErrorImplementingInterface(t *testing.T) {
	t.Parallel()

	err := stubErr{msg: "internal", userMsg: "please try again"}
	if got := ae.UserMessage(err); got != "please try again" {
		t.Errorf("UserMessage(stubErr) = %q, want %q", got, "please try again")
	}
}

func TestUserMessage_AeBuilderSetsUserMessage(t *testing.T) {
	t.Parallel()

	err := ae.New().UserMsg("internal", "user-safe")
	if got := ae.UserMessage(err); got != "user-safe" {
		t.Errorf("UserMessage on builder = %q, want %q", got, "user-safe")
	}
	if got := ae.Message(err); got != "internal" {
		t.Errorf("Message on builder = %q, want %q", got, "internal")
	}
}
