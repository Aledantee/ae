package ae_test

import (
	"errors"
	"reflect"
	"slices"
	"strings"
	"testing"

	"go.aledante.io/ae"
)

func TestAe_ErrorReturnsMessageOnly(t *testing.T) {
	t.Parallel()

	err := ae.New().Msg("disk full")
	if got := err.Error(); got != "disk full" {
		t.Errorf("Error() = %q, want %q", got, "disk full")
	}
}

func TestAe_ErrorIncludesSingleCause(t *testing.T) {
	t.Parallel()

	cause := errors.New("nvme offline")
	err := ae.New().Cause(cause).Msg("disk full")
	want := "disk full: nvme offline"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestAe_ErrorBracketsMultipleCauses(t *testing.T) {
	t.Parallel()

	c1 := errors.New("a")
	c2 := errors.New("b")
	c3 := errors.New("c")
	err := ae.New().Cause(c1, c2, c3).Msg("multi")
	want := "multi: [a; b; c]"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestAe_UnwrapReturnsCauses(t *testing.T) {
	t.Parallel()

	c1 := errors.New("one")
	c2 := errors.New("two")
	err := ae.New().Cause(c1, c2).Msg("outer")

	// errors.Is walks Unwrap() []error.
	if !errors.Is(err, c1) {
		t.Errorf("errors.Is did not find c1 through Unwrap")
	}
	if !errors.Is(err, c2) {
		t.Errorf("errors.Is did not find c2 through Unwrap")
	}
}

func TestAe_ErrorMessageAccessor(t *testing.T) {
	t.Parallel()

	err := ae.New().Msg("abc")
	em, ok := err.(ae.ErrorMessage)
	if !ok {
		t.Fatal("*ae.Ae does not implement ErrorMessage")
	}
	if em.ErrorMessage() != "abc" {
		t.Errorf("ErrorMessage = %q, want %q", em.ErrorMessage(), "abc")
	}
}

func TestAe_AllAccessorsReflectBuilder(t *testing.T) {
	t.Parallel()

	err := ae.New().
		Code("E").
		Hint("h").
		ExitCode(4).
		TraceId("t").
		SpanId("s").
		Tag("one").
		Attr("k", "v").
		UserMsg("internal", "friendly")

	if ae.Message(err) != "internal" {
		t.Errorf("Message = %q, want internal", ae.Message(err))
	}
	if ae.UserMessage(err) != "friendly" {
		t.Errorf("UserMessage = %q, want friendly", ae.UserMessage(err))
	}
	if ae.Code(err) != "E" {
		t.Errorf("Code = %q, want E", ae.Code(err))
	}
	if ae.Hint(err) != "h" {
		t.Errorf("Hint = %q, want h", ae.Hint(err))
	}
	if ae.ExitCode(err) != 4 {
		t.Errorf("ExitCode = %d, want 4", ae.ExitCode(err))
	}
	if ae.TraceId(err) != "t" {
		t.Errorf("TraceId = %q, want t", ae.TraceId(err))
	}
	if ae.SpanId(err) != "s" {
		t.Errorf("SpanId = %q, want s", ae.SpanId(err))
	}

	tags := ae.Tags(err)
	if !slices.Contains(tags, "one") {
		t.Errorf("Tags = %v, want to contain 'one'", tags)
	}

	attrs := ae.Attributes(err)
	if attrs["k"] != "v" {
		t.Errorf("Attributes k = %v, want 'v'", attrs["k"])
	}
}

func TestAe_ErrorTagsReturnsCopy(t *testing.T) {
	t.Parallel()

	err := ae.New().Tag("one").Msg("x")
	a, ok := err.(ae.ErrorTags)
	if !ok {
		t.Fatal("error does not implement ErrorTags")
	}

	first := a.ErrorTags()
	first = append(first, "mutated")
	_ = first

	second := a.ErrorTags()
	if slices.Contains(second, "mutated") {
		t.Errorf("ErrorTags leaked mutation: %v", second)
	}
}

func TestAe_ErrorAttributesReturnsCopy(t *testing.T) {
	t.Parallel()

	err := ae.New().Attr("k", "v").Msg("x")
	a, ok := err.(ae.ErrorAttributes)
	if !ok {
		t.Fatal("error does not implement ErrorAttributes")
	}

	first := a.ErrorAttributes()
	first["k"] = "mutated"
	first["extra"] = true

	second := a.ErrorAttributes()
	if second["k"] != "v" {
		t.Errorf("ErrorAttributes leaked k mutation: %v", second)
	}
	if _, present := second["extra"]; present {
		t.Errorf("ErrorAttributes leaked new-key mutation: %v", second)
	}
}

func TestAe_ErrorCausesReturnsCopy(t *testing.T) {
	t.Parallel()

	orig := errors.New("orig")
	err := ae.New().Cause(orig).Msg("x")
	a, ok := err.(ae.ErrorCauses)
	if !ok {
		t.Fatal("error does not implement ErrorCauses")
	}

	first := a.ErrorCauses()
	if len(first) == 0 {
		t.Fatal("expected at least one cause")
	}
	first[0] = errors.New("tampered")

	second := a.ErrorCauses()
	if second[0].Error() != "orig" {
		t.Errorf("ErrorCauses leaked slice mutation: got %v, want orig", second[0])
	}
}

func TestAe_ErrorRelatedReturnsCopy(t *testing.T) {
	t.Parallel()

	r := errors.New("r1")
	err := ae.New().Related(r).Msg("x")
	a, ok := err.(ae.ErrorRelated)
	if !ok {
		t.Fatal("error does not implement ErrorRelated")
	}

	first := a.ErrorRelated()
	if len(first) == 0 {
		t.Fatal("expected at least one related")
	}
	first[0] = errors.New("tampered")

	second := a.ErrorRelated()
	if second[0].Error() != "r1" {
		t.Errorf("ErrorRelated leaked slice mutation: got %v, want r1", second[0])
	}
}

func TestAe_PrintsReturnsNonEmptyString(t *testing.T) {
	t.Parallel()

	err := ae.New().Code("X").Msg("m")
	// Ae's own Prints method defers to the printer; just verify it runs and
	// yields the documented kind of output.
	aeErr, ok := err.(interface {
		Prints(opts ...ae.PrinterOption) string
	})
	if !ok {
		t.Fatal("*ae.Ae does not expose Prints")
	}
	out := aeErr.Prints(ae.NoPrintColors())
	if !strings.Contains(out, "m") {
		t.Errorf("Prints output = %q, want to contain message %q", out, "m")
	}
}

func TestAe_ImplementsAllInterfaces(t *testing.T) {
	t.Parallel()

	// The documented contract of *Ae is that it implements every ErrorX
	// interface in the package. Pin that here so a future refactor that drops
	// one surfaces immediately.
	var err error = ae.New().Msg("x")

	checks := []struct {
		name string
		ok   bool
	}{
		{"ErrorMessage", implements[ae.ErrorMessage](err)},
		{"ErrorUserMessage", implements[ae.ErrorUserMessage](err)},
		{"ErrorHint", implements[ae.ErrorHint](err)},
		{"ErrorCode", implements[ae.ErrorCode](err)},
		{"ErrorExitCode", implements[ae.ErrorExitCode](err)},
		{"ErrorTraceId", implements[ae.ErrorTraceId](err)},
		{"ErrorSpanId", implements[ae.ErrorSpanId](err)},
		{"ErrorTags", implements[ae.ErrorTags](err)},
		{"ErrorAttributes", implements[ae.ErrorAttributes](err)},
		{"ErrorCauses", implements[ae.ErrorCauses](err)},
		{"ErrorRelated", implements[ae.ErrorRelated](err)},
		{"ErrorStacks", implements[ae.ErrorStacks](err)},
		{"ErrorTimestamp", implements[ae.ErrorTimestamp](err)},
	}
	for _, c := range checks {
		if !c.ok {
			t.Errorf("*ae.Ae does not implement %s", c.name)
		}
	}
	// Silence reflect unused import warnings when checks is refactored.
	_ = reflect.TypeOf(err)
}

func implements[T any](v any) bool {
	_, ok := v.(T)
	return ok
}
