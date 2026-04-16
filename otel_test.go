package ae_test

import (
	"context"
	"errors"
	"testing"

	"go.opentelemetry.io/otel/attribute"

	"go.aledante.io/ae"
)

func TestTraceId_NilError(t *testing.T) {
	t.Parallel()

	if got := ae.TraceId(nil); got != "" {
		t.Errorf("TraceId(nil) = %q, want empty string", got)
	}
}

func TestTraceId_ErrorWithoutInterface(t *testing.T) {
	t.Parallel()

	if got := ae.TraceId(errors.New("plain")); got != "" {
		t.Errorf("TraceId(plainErr) = %q, want empty string", got)
	}
}

func TestTraceId_ErrorImplementingInterface(t *testing.T) {
	t.Parallel()

	err := stubErr{msg: "x", traceId: "abc123"}
	if got := ae.TraceId(err); got != "abc123" {
		t.Errorf("TraceId(stubErr) = %q, want %q", got, "abc123")
	}
}

func TestSpanId_NilError(t *testing.T) {
	t.Parallel()

	if got := ae.SpanId(nil); got != "" {
		t.Errorf("SpanId(nil) = %q, want empty string", got)
	}
}

func TestSpanId_ErrorWithoutInterface(t *testing.T) {
	t.Parallel()

	if got := ae.SpanId(errors.New("plain")); got != "" {
		t.Errorf("SpanId(plainErr) = %q, want empty string", got)
	}
}

func TestSpanId_ErrorImplementingInterface(t *testing.T) {
	t.Parallel()

	err := stubErr{msg: "x", spanId: "span-9"}
	if got := ae.SpanId(err); got != "span-9" {
		t.Errorf("SpanId(stubErr) = %q, want %q", got, "span-9")
	}
}

func TestBuilder_TraceIdAndSpanId(t *testing.T) {
	t.Parallel()

	err := ae.New().TraceId("tid").SpanId("sid").Msg("x")
	if got := ae.TraceId(err); got != "tid" {
		t.Errorf("TraceId = %q, want %q", got, "tid")
	}
	if got := ae.SpanId(err); got != "sid" {
		t.Errorf("SpanId = %q, want %q", got, "sid")
	}
}

// TestWithOtelAttribute_Roundtrip asserts that an attribute added through
// WithOtelAttribute can be read back via AttributesFromContext. The helpers
// forward to WithAttribute, so this test is blocked by the attributes.go key
// mismatch and is expected to fail until that is fixed.
func TestWithOtelAttribute_Roundtrip(t *testing.T) {
	t.Parallel()

	ctx := ae.WithOtelAttribute(context.Background(), attribute.String("service", "auth"))
	got := ae.AttributesFromContext(ctx)
	if got["service"] == nil {
		t.Errorf("AttributesFromContext after WithOtelAttribute = %v, want service attribute present", got)
	}
}

func TestWithOtelAttributeSet_Roundtrip(t *testing.T) {
	t.Parallel()

	set := attribute.NewSet(attribute.String("k", "v"), attribute.Int("n", 1))
	ctx := ae.WithOtelAttributeSet(context.Background(), set)
	got := ae.AttributesFromContext(ctx)
	if got["k"] == nil || got["n"] == nil {
		t.Errorf("AttributesFromContext after WithOtelAttributeSet = %v, want k and n present", got)
	}
}
