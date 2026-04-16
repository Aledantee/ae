package ae_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"go.aledante.io/ae"
)

func TestStacks_NilError(t *testing.T) {
	t.Parallel()

	if got := ae.Stacks(nil); got != nil {
		t.Errorf("Stacks(nil) = %v, want nil", got)
	}
}

func TestStacks_ErrorWithoutInterface(t *testing.T) {
	t.Parallel()

	if got := ae.Stacks(errors.New("plain")); got != nil {
		t.Errorf("Stacks(plainErr) = %v, want nil", got)
	}
}

func TestStacks_ErrorImplementingInterface(t *testing.T) {
	t.Parallel()

	sample := []*ae.Stack{{ID: 42, State: "running"}}
	err := stubErr{msg: "x", stacks: sample}
	got := ae.Stacks(err)

	if !reflect.DeepEqual(got, sample) {
		t.Errorf("Stacks = %v, want %v", got, sample)
	}
}

func TestBuilder_StackCapturesAtLeastOneGoroutine(t *testing.T) {
	t.Parallel()

	err := ae.New().Stack().Msg("with-stack")
	stacks := ae.Stacks(err)
	if len(stacks) == 0 {
		t.Fatal("Stack() produced no stacks")
	}
}

func TestBuilder_StackDropsDocumentedHelpers(t *testing.T) {
	t.Parallel()

	err := ae.New().Stack().Msg("with-stack")
	stacks := ae.Stacks(err)
	if len(stacks) == 0 {
		t.Fatal("Stack() produced no stacks")
	}

	// The package documents that ae.newStack, ae.Builder.Stack, and debug.Stack
	// are filtered from the captured frames. Assert none appear.
	dropped := []string{"ae.newStack", "ae.Builder.Stack", "debug.Stack"}
	for _, stack := range stacks {
		for _, f := range stack.Frames {
			for _, d := range dropped {
				if strings.HasSuffix(f.Func, d) {
					t.Errorf("frame %q leaked through the drop filter %q", f.Func, d)
				}
			}
		}
	}
}

func TestStackFrame_FieldsExported(t *testing.T) {
	t.Parallel()

	// The Stack and StackFrame structs are part of the public API; their
	// fields are documented. Pin the public shape so a refactor can't
	// silently rename or unexport one.
	f := ae.StackFrame{Func: "f", File: "x.go", Line: 7}
	if f.Func != "f" || f.File != "x.go" || f.Line != 7 {
		t.Errorf("StackFrame fields not accessible as expected: %+v", f)
	}

	s := ae.Stack{
		ID: 1, State: "running", Locked: true,
		Frames:    []*ae.StackFrame{&f},
		CreatedBy: &f,
		Ancestor:  nil,
	}
	if s.ID != 1 || s.State != "running" || !s.Locked {
		t.Errorf("Stack field access not as expected: %+v", s)
	}
	if len(s.Frames) != 1 || s.Frames[0].Func != "f" {
		t.Errorf("Stack.Frames not as expected: %+v", s.Frames)
	}
}
