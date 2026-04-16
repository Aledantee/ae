package ae_test

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"go.aledante.io/ae"
)

// buildRichErr builds an error touching every documented field that a printer
// might render. Reused across printer tests.
func buildRichErr(t *testing.T) error {
	t.Helper()
	inner := errors.New("root cause")
	return ae.New().
		Code("E_AUTH").
		ExitCode(2).
		Hint("try again").
		Tag("network").
		Attr("attempt", 3).
		TraceId("trace-x").
		SpanId("span-y").
		Cause(inner).
		Related(errors.New("side-issue")).
		Msg("failed")
}

func TestNewPrinter_DefaultsIncludeVerboseText(t *testing.T) {
	t.Parallel()

	err := buildRichErr(t)
	// Strip colors to keep assertions stable; defaults are verbose and text.
	out := ae.NewPrinter(ae.NoPrintColors()).Prints(err)

	wants := []string{
		"failed",      // message
		"E_AUTH",      // code
		"network",     // tag (verbose default includes tags)
		"try again",   // hint
		"attempt",     // attribute key
		"root cause",  // cause message
	}
	for _, w := range wants {
		if !strings.Contains(out, w) {
			t.Errorf("output missing %q\n--- output ---\n%s", w, out)
		}
	}
}

func TestNewPrinter_JSONProducesValidJSONWithDocumentedKeys(t *testing.T) {
	t.Parallel()

	err := buildRichErr(t)
	out := ae.NewPrinter(ae.PrintJSON()).Prints(err)

	var got map[string]any
	if decodeErr := json.Unmarshal([]byte(out), &got); decodeErr != nil {
		t.Fatalf("JSON output did not parse: %v\n%s", decodeErr, out)
	}

	// jsonError struct tags document these keys. Only non-empty values are
	// emitted (omitempty), so we assert each is present and has the right
	// type/value for the rich fixture.
	wants := map[string]any{
		"message":   "failed",
		"hint":      "try again",
		"code":      "E_AUTH",
		"trace_id":  "trace-x",
		"span_id":   "span-y",
	}
	for k, v := range wants {
		if got[k] != v {
			t.Errorf("JSON[%q] = %v, want %v", k, got[k], v)
		}
	}

	// exit_code is int-typed in the struct but decodes as float64 through
	// encoding/json's default.
	if ec, ok := got["exit_code"].(float64); !ok || ec != 2 {
		t.Errorf("JSON[exit_code] = %v, want 2", got["exit_code"])
	}

	if _, ok := got["tags"].([]any); !ok {
		t.Errorf("JSON[tags] = %v, want []any", got["tags"])
	}
	if _, ok := got["attrs"].(map[string]any); !ok {
		t.Errorf("JSON[attrs] = %v, want map[string]any", got["attrs"])
	}
	if _, ok := got["causes"].([]any); !ok {
		t.Errorf("JSON[causes] = %v, want []any", got["causes"])
	}
	if _, ok := got["related"].([]any); !ok {
		t.Errorf("JSON[related] = %v, want []any", got["related"])
	}
}

func TestPrinter_PrintDepthZeroSuppressesCauses(t *testing.T) {
	t.Parallel()

	err := buildRichErr(t)
	out := ae.NewPrinter(ae.NoPrintColors(), ae.PrintDepth(0)).Prints(err)

	if strings.Contains(out, "root cause") {
		t.Errorf("PrintDepth(0) emitted cause text:\n%s", out)
	}
	if !strings.Contains(out, "failed") {
		t.Errorf("PrintDepth(0) lost outer message:\n%s", out)
	}
}

func TestPrinter_PrintDepthOneIncludesImmediateCause(t *testing.T) {
	t.Parallel()

	inner2 := errors.New("inner-2")
	inner1 := ae.New().Cause(inner2).Msg("inner-1")
	outer := ae.New().Cause(inner1).Msg("outer")

	out := ae.NewPrinter(ae.NoPrintColors(), ae.PrintDepth(1)).Prints(outer)

	if !strings.Contains(out, "inner-1") {
		t.Errorf("PrintDepth(1) dropped depth-1 cause:\n%s", out)
	}
	if strings.Contains(out, "inner-2") {
		t.Errorf("PrintDepth(1) leaked depth-2 cause:\n%s", out)
	}
}

// TestPrinter_NoPrintCausesSuppressesCauseBlock asserts the docstring of
// NoPrintCauses: "disables inclusion of error causes in the output."
// NewPrinter currently unconditionally appends PrintCauses() after all
// user-provided options, so NoPrintCauses is silently neutralized. This test
// is expected to fail until the forced append is removed.
func TestPrinter_NoPrintCausesSuppressesCauseBlock(t *testing.T) {
	t.Parallel()

	err := buildRichErr(t)
	out := ae.NewPrinter(ae.NoPrintColors(), ae.NoPrintCauses()).Prints(err)
	if strings.Contains(out, "root cause") {
		t.Errorf("NoPrintCauses still emitted cause:\n%s", out)
	}
}

func TestPrinter_NoPrintStacksSuppressesStackSection(t *testing.T) {
	t.Parallel()

	err := ae.New().Stack().Msg("with-stack")
	out := ae.NewPrinter(ae.NoPrintColors(), ae.NoPrintStacks()).Prints(err)

	// The stack block header renders as "  stack" on its own line.
	if strings.Contains(out, "\n  stack") {
		t.Errorf("NoPrintStacks still included stack header:\n%s", out)
	}
	if strings.Contains(out, "goroutine ") {
		t.Errorf("NoPrintStacks still emitted a goroutine frame:\n%s", out)
	}
}

func TestPrinter_PrintIndentClampsToOne(t *testing.T) {
	t.Parallel()

	err := buildRichErr(t)
	// PrintIndent(0) clamps to 1. The assertion: attributes still render and
	// the output does not panic.
	out := ae.NewPrinter(ae.NoPrintColors(), ae.PrintIndent(0)).Prints(err)
	if !strings.Contains(out, "attempt") {
		t.Errorf("PrintIndent(0) dropped attributes:\n%s", out)
	}
}

// TestPrinter_PrintVerboseIncludesAllSections asserts the PrintVerbose
// docstring: "enables all available output fields ... including ... trace IDs,
// span IDs ...". The current implementation wires PrintOtel to toggle only the
// trace id field, never the span id, so the span-y assertion fails until the
// printer is fixed.
func TestPrinter_PrintVerboseIncludesAllSections(t *testing.T) {
	t.Parallel()

	err := buildRichErr(t)
	out := ae.NewPrinter(ae.NoPrintColors(), ae.PrintVerbose()).Prints(err)
	for _, w := range []string{"trace-x", "span-y", "network", "attempt", "root cause"} {
		if !strings.Contains(out, w) {
			t.Errorf("PrintVerbose missing %q:\n%s", w, out)
		}
	}
}

func TestPrinter_PrintCompactIncludesDocumentedCompactFields(t *testing.T) {
	t.Parallel()

	err := ae.New().
		Code("C").
		Hint("h").
		Tag("t").
		Attr("k", "v").
		Msg("m")

	// PrintCompact enables hint, code, exit code, attributes, tags, causes,
	// related. It does not say anything about stacks, so this test only
	// asserts what PrintCompact documents: the "minimal set of commonly
	// useful output fields" includes each of these.
	out := ae.NewPrinter(ae.NoPrintColors(), ae.PrintCompact()).Prints(err)
	for _, w := range []string{"m", "C", "h", "t"} {
		if !strings.Contains(out, w) {
			t.Errorf("PrintCompact missing documented field %q:\n%s", w, out)
		}
	}
}

func TestPrinter_PrintColorsInjectsAnsiEscapes(t *testing.T) {
	t.Parallel()

	err := ae.New().Msg("hello")
	out := ae.NewPrinter(ae.PrintColors()).Prints(err)

	// ANSI SGR introducer is ESC[ — we assert only that *some* SGR sequence
	// appears when colors are on. (The fatih/color library no-ops when
	// stdout isn't a TTY, so this test may be skipped at runtime.)
	if strings.Contains(out, "\x1b[") {
		return
	}
	t.Skip("fatih/color suppresses ANSI escapes in non-TTY test environment; PrintColors code path is still exercised")
}

func TestPrinter_NoPrintColorsRemovesAnsiEscapes(t *testing.T) {
	t.Parallel()

	err := ae.New().Code("C").Msg("x")
	out := ae.NewPrinter(ae.NoPrintColors()).Prints(err)
	if strings.Contains(out, "\x1b[") {
		t.Errorf("NoPrintColors still emitted ANSI SGR: %q", out)
	}
}

func TestPrinter_NoPrintJSONProducesNonJSONOutput(t *testing.T) {
	t.Parallel()

	err := ae.New().Msg("plain")
	out := ae.NewPrinter(ae.NoPrintColors(), ae.NoPrintJSON()).Prints(err)
	if strings.HasPrefix(strings.TrimSpace(out), "{") {
		t.Errorf("NoPrintJSON still produced JSON-looking output: %q", out)
	}
}

// TestPrinter_TextRendersAttributeValue guards the attribute-value rendering
// path: the integer 3 must appear in the output and the literal format verb
// "%v" must not leak. Earlier regressions left "attempt: %v" in no-color mode.
func TestPrinter_TextRendersAttributeValue(t *testing.T) {
	t.Parallel()

	err := ae.New().Attr("attempt", 3).Msg("m")
	out := ae.NewPrinter(ae.NoPrintColors()).Prints(err)

	if !strings.Contains(out, "attempt") || !strings.Contains(out, "3") {
		t.Errorf("text attribute line missing rendered key/value:\n%s", out)
	}
	if strings.Contains(out, "%v") {
		t.Errorf("text attribute line contains literal '%%v' format verb:\n%s", out)
	}
}

func TestGlobalPrint_DoesNotPanic(t *testing.T) {
	t.Parallel()

	// Package-level ae.Print writes to stdout; just make sure it can be
	// called with a representative error without panicking.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ae.Print panicked: %v", r)
		}
	}()
	ae.Print(buildRichErr(t), ae.NoPrintColors())
}

func TestPrintFrameFilters_CustomFilterDropsMatchingFrames(t *testing.T) {
	t.Parallel()

	err := ae.New().Stack().Msg("with-stack")

	// Ask the printer to additionally drop any frame whose function contains
	// "testing." — net effect for a test process: strip the go test runtime
	// frames so only this test's own invocation path survives.
	filter := func(f *ae.StackFrame) bool {
		return f != nil && strings.Contains(f.Func, "testing.")
	}

	out := ae.NewPrinter(ae.NoPrintColors(), ae.PrintFrameFilters(filter)).Prints(err)

	if strings.Contains(out, "testing.") {
		t.Errorf("custom frame filter did not drop 'testing.' frames:\n%s", out)
	}
}

func TestPrintFrameFilters_GoroutineWithAllFramesFilteredIsOmitted(t *testing.T) {
	t.Parallel()

	err := ae.New().Stack().Msg("with-stack")

	// A filter that drops every frame leaves no frames to print. The
	// goroutine header must be omitted too — otherwise the output would
	// carry a stranded "goroutine N (state)" line with nothing below it.
	dropAll := func(f *ae.StackFrame) bool { return true }
	out := ae.NewPrinter(ae.NoPrintColors(), ae.PrintFrameFilters(dropAll)).Prints(err)

	if strings.Contains(out, "goroutine ") {
		t.Errorf("goroutine header leaked despite all frames filtered out:\n%s", out)
	}
}

func TestPrinter_FprintWritesToArbitraryWriter(t *testing.T) {
	t.Parallel()

	var buf strings.Builder
	ae.NewPrinter(ae.NoPrintColors()).Fprint(&buf, ae.New().Msg("hello"))

	got := buf.String()
	if !strings.Contains(got, "[ERROR]") || !strings.Contains(got, "hello") {
		t.Errorf("Fprint output missing expected substrings: %q", got)
	}
	if !strings.HasSuffix(got, "\n") {
		t.Errorf("Fprint output not newline-terminated: %q", got)
	}
}

func TestPrinter_PrintTraceIdAndPrintSpanIdIndependent(t *testing.T) {
	t.Parallel()

	err := ae.New().TraceId("tid").SpanId("sid").Msg("x")

	// Trace only — span should be absent.
	out := ae.NewPrinter(
		ae.NoPrintColors(),
		ae.NoPrintOtel(), // reset both
		ae.PrintTraceId(),
	).Prints(err)
	if !strings.Contains(out, "tid") {
		t.Errorf("PrintTraceId alone dropped trace id:\n%s", out)
	}
	if strings.Contains(out, "sid") {
		t.Errorf("PrintTraceId alone emitted span id:\n%s", out)
	}

	// Span only — trace should be absent.
	out = ae.NewPrinter(
		ae.NoPrintColors(),
		ae.NoPrintOtel(),
		ae.PrintSpanId(),
	).Prints(err)
	if !strings.Contains(out, "sid") {
		t.Errorf("PrintSpanId alone dropped span id:\n%s", out)
	}
	if strings.Contains(out, "tid") {
		t.Errorf("PrintSpanId alone emitted trace id:\n%s", out)
	}
}
