package ae_test

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"
	"time"

	"go.opentelemetry.io/otel/trace"

	"go.aledante.io/ae"
)

func TestNew_ReturnsUsableBuilder(t *testing.T) {
	t.Parallel()

	err := ae.New().Msg("ok")
	if err == nil {
		t.Fatal("New().Msg(...) returned nil")
	}
	if err.Error() != "ok" {
		t.Errorf("Error() = %q, want %q", err.Error(), "ok")
	}
}

func TestNewC_ShorthandForContext(t *testing.T) {
	t.Parallel()

	ctx := ae.WithTagsValue(context.Background(), "ctx-tag")
	err := ae.NewC(ctx).Msg("x")

	if !slices.Contains(ae.Tags(err), "ctx-tag") {
		t.Errorf("Tags = %v, want to contain 'ctx-tag' via NewC", ae.Tags(err))
	}
}

func TestFrom_NilErrorReturnsEmptyBuilder(t *testing.T) {
	t.Parallel()

	err := ae.From(nil).Msg("fresh")
	if err.Error() != "fresh" {
		t.Errorf("Error() = %q, want %q", err.Error(), "fresh")
	}
}

func TestFrom_AeClonesAllFields(t *testing.T) {
	t.Parallel()

	original := ae.New().
		Code("C").
		Hint("h").
		ExitCode(3).
		TraceId("t").
		SpanId("s").
		Tag("one").
		Attr("k", "v").
		Cause(errors.New("inner")).
		UserMsg("internal-orig", "user-orig")

	derived := ae.From(original).Msg("derived")

	if ae.Message(derived) != "derived" {
		t.Errorf("Message = %q, want 'derived'", ae.Message(derived))
	}
	if ae.UserMessage(derived) != "user-orig" {
		t.Errorf("UserMessage = %q, want 'user-orig'", ae.UserMessage(derived))
	}
	if ae.Code(derived) != "C" {
		t.Errorf("Code = %q, want 'C'", ae.Code(derived))
	}
	if ae.Hint(derived) != "h" {
		t.Errorf("Hint = %q, want 'h'", ae.Hint(derived))
	}
	if ae.ExitCode(derived) != 3 {
		t.Errorf("ExitCode = %d, want 3", ae.ExitCode(derived))
	}
	if ae.TraceId(derived) != "t" {
		t.Errorf("TraceId = %q, want 't'", ae.TraceId(derived))
	}
	if ae.SpanId(derived) != "s" {
		t.Errorf("SpanId = %q, want 's'", ae.SpanId(derived))
	}
	if !slices.Contains(ae.Tags(derived), "one") {
		t.Errorf("Tags = %v, want to contain 'one'", ae.Tags(derived))
	}
	if ae.Attributes(derived)["k"] != "v" {
		t.Errorf("Attributes k = %v, want 'v'", ae.Attributes(derived)["k"])
	}
	if len(ae.Causes(derived)) != 1 {
		t.Errorf("Causes = %v, want one", ae.Causes(derived))
	}
}

func TestFrom_ArbitraryErrorPullsInterfaceValues(t *testing.T) {
	t.Parallel()

	src := stubErr{
		msg:      "m",
		userMsg:  "u",
		code:     "X",
		exitCode: 2,
		hint:     "hh",
		traceId:  "tt",
		spanId:   "ss",
		tags:     []string{"alpha"},
		attrs:    map[string]any{"k": "v"},
	}
	err := ae.From(src).Msg("wrapped")

	if ae.Message(err) != "wrapped" {
		t.Errorf("Message = %q, want 'wrapped'", ae.Message(err))
	}
	if ae.UserMessage(err) != "u" {
		t.Errorf("UserMessage = %q, want 'u'", ae.UserMessage(err))
	}
	if ae.Code(err) != "X" {
		t.Errorf("Code = %q, want 'X'", ae.Code(err))
	}
	if ae.ExitCode(err) != 2 {
		t.Errorf("ExitCode = %d, want 2", ae.ExitCode(err))
	}
	if ae.Hint(err) != "hh" {
		t.Errorf("Hint = %q, want 'hh'", ae.Hint(err))
	}
	if ae.TraceId(err) != "tt" {
		t.Errorf("TraceId = %q, want 'tt'", ae.TraceId(err))
	}
	if ae.SpanId(err) != "ss" {
		t.Errorf("SpanId = %q, want 'ss'", ae.SpanId(err))
	}
	if !slices.Contains(ae.Tags(err), "alpha") {
		t.Errorf("Tags = %v, want to contain 'alpha'", ae.Tags(err))
	}
	if ae.Attributes(err)["k"] != "v" {
		t.Errorf("Attributes k = %v, want 'v'", ae.Attributes(err)["k"])
	}
}

func TestFromC_CombinesErrorAndContext(t *testing.T) {
	t.Parallel()

	ctx := ae.WithTagsValue(context.Background(), "from-ctx")
	src := stubErr{msg: "m", tags: []string{"from-err"}}

	err := ae.FromC(ctx, src).Msg("combined")
	tags := ae.Tags(err)

	if !slices.Contains(tags, "from-ctx") {
		t.Errorf("Tags = %v, want to contain 'from-ctx'", tags)
	}
	if !slices.Contains(tags, "from-err") {
		t.Errorf("Tags = %v, want to contain 'from-err'", tags)
	}
}

func TestBuilder_MsgIsTerminalAndReturnsError(t *testing.T) {
	t.Parallel()

	err := ae.New().Msg("boom")
	if err == nil {
		t.Fatal("Msg returned nil")
	}
	if err.Error() != "boom" {
		t.Errorf("Error() = %q, want %q", err.Error(), "boom")
	}
}

func TestBuilder_MsgfFormats(t *testing.T) {
	t.Parallel()

	err := ae.New().Msgf("code=%d name=%s", 42, "x")
	want := "code=42 name=x"
	if err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
}

func TestBuilder_UserMsgSetsBothMessages(t *testing.T) {
	t.Parallel()

	err := ae.New().UserMsg("internal", "external")
	if ae.Message(err) != "internal" {
		t.Errorf("Message = %q, want 'internal'", ae.Message(err))
	}
	if ae.UserMessage(err) != "external" {
		t.Errorf("UserMessage = %q, want 'external'", ae.UserMessage(err))
	}
}

func TestBuilder_CausesFiltersNil(t *testing.T) {
	t.Parallel()

	c := errors.New("real")
	err := ae.New().Cause(nil, c, nil).Msg("x")

	causes := ae.Causes(err)
	if len(causes) != 1 || causes[0] != c {
		t.Errorf("Causes = %v, want [real] (nil filtered)", causes)
	}
}

func TestBuilder_RelatedFiltersNil(t *testing.T) {
	t.Parallel()

	r := errors.New("real")
	err := ae.New().Related(nil, r, nil).Msg("x")

	rel := ae.Related(err)
	if len(rel) != 1 || rel[0] != r {
		t.Errorf("Related = %v, want [real] (nil filtered)", rel)
	}
}

func TestBuilder_CauseUnwrapExpandsMultiError(t *testing.T) {
	t.Parallel()

	c1 := errors.New("a")
	c2 := errors.New("b")
	wrapped := multiUnwrapErr{msg: "wrap", errs: []error{c1, c2}}

	err := ae.New().CauseUnwrap(wrapped).Msg("x")
	causes := ae.Causes(err)

	if len(causes) != 2 {
		t.Fatalf("Causes = %v, want 2 expanded entries", causes)
	}
	if !errors.Is(err, c1) || !errors.Is(err, c2) {
		t.Errorf("errors.Is did not find expanded causes")
	}
}

func TestBuilder_CauseUnwrapPreservesRegularError(t *testing.T) {
	t.Parallel()

	c := errors.New("plain")
	err := ae.New().CauseUnwrap(c).Msg("x")

	causes := ae.Causes(err)
	if len(causes) != 1 || causes[0] != c {
		t.Errorf("Causes = %v, want [plain]", causes)
	}
}

func TestBuilder_RelatedUnwrapExpandsMultiError(t *testing.T) {
	t.Parallel()

	r1 := errors.New("x")
	r2 := errors.New("y")
	wrapped := multiUnwrapErr{msg: "wrap", errs: []error{r1, r2}}

	err := ae.New().RelatedUnwrap(wrapped).Msg("x")
	rel := ae.Related(err)

	if len(rel) != 2 {
		t.Fatalf("Related = %v, want 2 expanded entries", rel)
	}
}

func TestBuilder_ContextPullsSpanAndTraceIds(t *testing.T) {
	t.Parallel()

	ctx := traceContextWith(t, "1234567890abcdef1234567890abcdef", "abcdef1234567890")
	err := ae.NewC(ctx).Msg("x")

	if ae.TraceId(err) == "" {
		t.Errorf("TraceId = empty, want populated from context span")
	}
	if ae.SpanId(err) == "" {
		t.Errorf("SpanId = empty, want populated from context span")
	}
}

func TestBuilder_ContextAddsProvidedKeysAsAttributes(t *testing.T) {
	t.Parallel()

	type reqIDKey struct{}
	ctx := context.WithValue(context.Background(), reqIDKey{}, "req-7")

	// String keys: the doc says the key may be string, fmt.Stringer, or other.
	ctx2 := context.WithValue(ctx, "trace_name", "root")

	err := ae.New().Context(ctx2, reqIDKey{}, "trace_name").Msg("x")
	attrs := ae.Attributes(err)

	if attrs["trace_name"] != "root" {
		t.Errorf("Attributes[trace_name] = %v, want 'root'", attrs["trace_name"])
	}
	// The non-string key is stringified via fmt.Sprintf %v. The exact format
	// of an empty struct{} key is implementation-defined; verify the value
	// was stored under some key mapping to req-7.
	found := false
	for _, v := range attrs {
		if v == "req-7" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Attributes = %v, want a value 'req-7' from the typed key", attrs)
	}
}

func TestBuilder_TimestampExplicit(t *testing.T) {
	t.Parallel()

	when := time.Date(2030, 6, 1, 0, 0, 0, 0, time.UTC)
	err := ae.New().Timestamp(when).Msg("x")
	if got := ae.Timestamp(err); !got.Equal(when) {
		t.Errorf("Timestamp = %v, want %v", got, when)
	}
}

func TestBuilder_CauseStringIncludedInError(t *testing.T) {
	t.Parallel()

	cause := errors.New("root cause")
	err := ae.New().Cause(cause).Msg("wrapping")
	if !strings.Contains(err.Error(), "root cause") {
		t.Errorf("Error() = %q, want to contain 'root cause'", err.Error())
	}
}

// traceContextWith returns a context carrying a valid OpenTelemetry SpanContext
// built from the given hex strings. Kept here so tests that need a real span
// context don't pull in a full tracer.
func traceContextWith(t *testing.T, traceHex, spanHex string) context.Context {
	t.Helper()

	tid, err := trace.TraceIDFromHex(traceHex)
	if err != nil {
		t.Fatalf("bad trace hex: %v", err)
	}
	sid, err := trace.SpanIDFromHex(spanHex)
	if err != nil {
		t.Fatalf("bad span hex: %v", err)
	}

	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: trace.FlagsSampled,
		Remote:     true,
	})
	return trace.ContextWithSpanContext(context.Background(), sc)
}
