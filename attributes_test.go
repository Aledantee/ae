package ae_test

import (
	"context"
	"errors"
	"testing"

	"go.aledante.io/ae"
)

func TestAttributes_NilError(t *testing.T) {
	t.Parallel()

	got := ae.Attributes(nil)
	if got == nil {
		t.Fatal("Attributes(nil) = nil, want empty non-nil map")
	}
	if len(got) != 0 {
		t.Errorf("Attributes(nil) = %v, want empty map", got)
	}
}

func TestAttributes_ErrorWithoutInterface(t *testing.T) {
	t.Parallel()

	got := ae.Attributes(errors.New("plain"))
	if got == nil {
		t.Fatal("Attributes(plainErr) = nil, want empty non-nil map")
	}
	if len(got) != 0 {
		t.Errorf("Attributes(plainErr) = %v, want empty map", got)
	}
}

func TestAttributes_ErrorImplementingInterface(t *testing.T) {
	t.Parallel()

	err := stubErr{msg: "x", attrs: map[string]any{"key": "val", "n": 42}}
	got := ae.Attributes(err)
	if got["key"] != "val" || got["n"] != 42 {
		t.Errorf("Attributes(stubErr) = %v, want map with key=val, n=42", got)
	}
}

func TestAttributes_AeBuilderAttrAndAttrs(t *testing.T) {
	t.Parallel()

	err := ae.New().
		Attr("a", 1).
		Attrs(map[string]any{"b": "two", "c": 3.14}).
		Msg("x")

	got := ae.Attributes(err)
	if got["a"] != 1 || got["b"] != "two" || got["c"] != 3.14 {
		t.Errorf("Attributes after Attr+Attrs = %v, want {a:1, b:two, c:3.14}", got)
	}
}

func TestAttributesFromContext_EmptyContext(t *testing.T) {
	t.Parallel()

	got := ae.AttributesFromContext(context.Background())
	if got == nil {
		t.Fatal("AttributesFromContext(bg) = nil, want empty non-nil map")
	}
	if len(got) != 0 {
		t.Errorf("AttributesFromContext(bg) = %v, want empty map", got)
	}
}

// TestAttributesFromContext_Roundtrip asserts the documented roundtrip between
// WithAttribute and AttributesFromContext. The current implementation reads
// tagKey{} but WithAttribute writes under attributesKey{}; expected to fail
// until the key mismatch is fixed.
func TestAttributesFromContext_Roundtrip(t *testing.T) {
	t.Parallel()

	ctx := ae.WithAttribute(context.Background(), "user_id", "u-42")
	got := ae.AttributesFromContext(ctx)
	if got["user_id"] != "u-42" {
		t.Errorf("AttributesFromContext after WithAttribute = %v, want user_id=u-42", got)
	}
}

// TestWithAttributes_OverwritesDuplicateKeys asserts documented overwrite
// semantics for WithAttributes. Same key-mismatch bug: expected to fail.
func TestWithAttributes_OverwritesDuplicateKeys(t *testing.T) {
	t.Parallel()

	ctx := ae.WithAttributes(context.Background(), map[string]any{"k": "first"})
	ctx = ae.WithAttributes(ctx, map[string]any{"k": "second"})

	got := ae.AttributesFromContext(ctx)
	if got["k"] != "second" {
		t.Errorf("AttributesFromContext overwrite = %v, want k=second", got)
	}
}

// TestBuilder_ContextPullsAttributesIntoError asserts that an error built with
// NewC(ctx) carries attributes attached to the context. Builder.Context calls
// AttributesFromContext, which is blocked by the same key-mismatch bug;
// expected to fail until that is fixed.
func TestBuilder_ContextPullsAttributesIntoError(t *testing.T) {
	t.Parallel()

	ctx := ae.WithAttribute(context.Background(), "request_id", "r-7")
	err := ae.NewC(ctx).Msg("x")
	got := ae.Attributes(err)
	if got["request_id"] != "r-7" {
		t.Errorf("Attributes after NewC = %v, want request_id=r-7", got)
	}
}
