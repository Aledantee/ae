package ae_test

// Note: Exit and PrintExit call os.Exit, which terminates the test process.
// They are thin wrappers over ExitCode (which is exhaustively tested in
// exitcode_test.go) plus os.Exit, and testing them directly would require a
// subprocess harness that is disproportionate to the two lines of logic they
// add. They are deliberately left untested here.

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"

	"go.aledante.io/ae"
)

func TestWrap_NilErrorReturnsNil(t *testing.T) {
	t.Parallel()

	if got := ae.Wrap("ctx", nil); got != nil {
		t.Errorf("Wrap(ctx, nil) = %v, want nil", got)
	}
}

func TestWrap_WrapsErrorWithMessage(t *testing.T) {
	t.Parallel()

	cause := errors.New("root")
	err := ae.Wrap("outer", cause)

	if err == nil {
		t.Fatal("Wrap returned nil for a non-nil cause")
	}
	if !strings.Contains(err.Error(), "outer") {
		t.Errorf("Error() = %q, want to contain 'outer'", err.Error())
	}
	if !errors.Is(err, cause) {
		t.Errorf("errors.Is did not find wrapped cause")
	}
	causes := ae.Causes(err)
	if len(causes) != 1 || causes[0] != cause {
		t.Errorf("Causes = %v, want [%v]", causes, cause)
	}
}

func TestWrapC_NilErrorReturnsNil(t *testing.T) {
	t.Parallel()

	if got := ae.WrapC(context.Background(), "ctx", nil); got != nil {
		t.Errorf("WrapC(..., nil) = %v, want nil", got)
	}
}

func TestWrapC_IncludesContextDerivedData(t *testing.T) {
	t.Parallel()

	ctx := ae.WithTagsValue(context.Background(), "from-ctx")
	err := ae.WrapC(ctx, "outer", errors.New("inner"))

	if !slices.Contains(ae.Tags(err), "from-ctx") {
		t.Errorf("Tags = %v, want to contain 'from-ctx'", ae.Tags(err))
	}
}

func TestWrapf_FormatsMessage(t *testing.T) {
	t.Parallel()

	inner := errors.New("inner")
	err := ae.Wrapf("count=%d name=%s", inner, 7, "x")
	if !strings.Contains(err.Error(), "count=7 name=x") {
		t.Errorf("Error() = %q, want to contain 'count=7 name=x'", err.Error())
	}
	if !errors.Is(err, inner) {
		t.Errorf("errors.Is could not find inner through Wrapf")
	}
}

func TestWrapCf_FormatsMessageAndPullsContext(t *testing.T) {
	t.Parallel()

	ctx := ae.WithTagsValue(context.Background(), "cf-tag")
	inner := errors.New("inner")
	err := ae.WrapCf(ctx, "level=%s", inner, "warn")

	if !strings.Contains(err.Error(), "level=warn") {
		t.Errorf("Error() = %q, want to contain 'level=warn'", err.Error())
	}
	if !slices.Contains(ae.Tags(err), "cf-tag") {
		t.Errorf("Tags = %v, want to contain 'cf-tag'", ae.Tags(err))
	}
}

func TestWrapMany_AllNilReturnsNil(t *testing.T) {
	t.Parallel()

	if got := ae.WrapMany("outer", nil, nil); got != nil {
		t.Errorf("WrapMany(outer, nil, nil) = %v, want nil", got)
	}
}

func TestWrapMany_NoArgsReturnsNil(t *testing.T) {
	t.Parallel()

	if got := ae.WrapMany("outer"); got != nil {
		t.Errorf("WrapMany(outer) = %v, want nil", got)
	}
}

func TestWrapMany_FiltersNilsAndWraps(t *testing.T) {
	t.Parallel()

	c1 := errors.New("a")
	c2 := errors.New("b")
	err := ae.WrapMany("outer", nil, c1, nil, c2)

	if err == nil {
		t.Fatal("WrapMany returned nil for non-empty input")
	}
	causes := ae.Causes(err)
	if len(causes) != 2 {
		t.Errorf("Causes = %v, want 2 (nils filtered)", causes)
	}
	if !errors.Is(err, c1) || !errors.Is(err, c2) {
		t.Errorf("errors.Is did not find both wrapped causes")
	}
}

func TestMsg_ProducesErrorWithMessage(t *testing.T) {
	t.Parallel()

	err := ae.Msg("boom")
	if err == nil || err.Error() != "boom" {
		t.Errorf("Msg = %v, want error with message 'boom'", err)
	}
}

func TestMsgC_AttachesContext(t *testing.T) {
	t.Parallel()

	ctx := ae.WithTagsValue(context.Background(), "mc-tag")
	err := ae.MsgC(ctx, "x")

	if err.Error() != "x" {
		t.Errorf("Error() = %q, want 'x'", err.Error())
	}
	if !slices.Contains(ae.Tags(err), "mc-tag") {
		t.Errorf("Tags = %v, want to contain 'mc-tag'", ae.Tags(err))
	}
}

func TestMsgf_Formats(t *testing.T) {
	t.Parallel()

	err := ae.Msgf("val=%d", 9)
	if err.Error() != "val=9" {
		t.Errorf("Error() = %q, want 'val=9'", err.Error())
	}
}

func TestMsgCf_FormatsAndAttachesContext(t *testing.T) {
	t.Parallel()

	ctx := ae.WithTagsValue(context.Background(), "mcf-tag")
	err := ae.MsgCf(ctx, "%s=%d", "count", 3)

	if err.Error() != "count=3" {
		t.Errorf("Error() = %q, want 'count=3'", err.Error())
	}
	if !slices.Contains(ae.Tags(err), "mcf-tag") {
		t.Errorf("Tags = %v, want to contain 'mcf-tag'", ae.Tags(err))
	}
}
