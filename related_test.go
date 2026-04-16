package ae_test

import (
	"errors"
	"testing"

	"go.aledante.io/ae"
)

func TestRelated_NilError(t *testing.T) {
	t.Parallel()

	if got := ae.Related(nil); got != nil {
		t.Errorf("Related(nil) = %v, want nil", got)
	}
}

func TestRelated_ErrorWithoutInterface(t *testing.T) {
	t.Parallel()

	if got := ae.Related(errors.New("plain")); got != nil {
		t.Errorf("Related(plainErr) = %v, want nil", got)
	}
}

func TestRelated_ErrorImplementingInterface(t *testing.T) {
	t.Parallel()

	r1 := errors.New("r1")
	r2 := errors.New("r2")
	err := stubErr{msg: "x", related: []error{r1, r2}}
	got := ae.Related(err)
	if len(got) != 2 || got[0] != r1 || got[1] != r2 {
		t.Errorf("Related(stubErr) = %v, want [r1 r2]", got)
	}
}
