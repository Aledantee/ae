package ae_test

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"testing"

	"go.aledante.io/ae"
)

func TestTags_NilError(t *testing.T) {
	t.Parallel()

	if got := ae.Tags(nil); got != nil {
		t.Errorf("Tags(nil) = %v, want nil", got)
	}
}

func TestTags_ErrorWithoutInterface(t *testing.T) {
	t.Parallel()

	if got := ae.Tags(errors.New("plain")); got != nil {
		t.Errorf("Tags(plainErr) = %v, want nil", got)
	}
}

func TestTags_ErrorImplementingInterface(t *testing.T) {
	t.Parallel()

	err := stubErr{msg: "x", tags: []string{"db", "network"}}
	got := ae.Tags(err)
	slices.Sort(got)
	want := []string{"db", "network"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Tags(stubErr) = %v, want %v", got, want)
	}
}

func TestTags_AeBuilderAddsAndDeduplicates(t *testing.T) {
	t.Parallel()

	err := ae.New().Tag("a").Tag("b").Tag("a").Tags("c", "b").Msg("x")
	got := ae.Tags(err)
	slices.Sort(got)
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Tags after dedupe = %v, want %v", got, want)
	}
}

func TestTagsFromContext_EmptyContext(t *testing.T) {
	t.Parallel()

	if got := ae.TagsFromContext(context.Background()); got != nil {
		t.Errorf("TagsFromContext(bg) = %v, want nil", got)
	}
}

func TestWithTagsValue_RoundtripsThroughContext(t *testing.T) {
	t.Parallel()

	ctx := ae.WithTagsValue(context.Background(), "alpha", "beta")
	got := ae.TagsFromContext(ctx)
	want := []string{"alpha", "beta"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("TagsFromContext = %v, want %v", got, want)
	}
}

func TestWithTagsValue_DeduplicatesAndSorts(t *testing.T) {
	t.Parallel()

	ctx := ae.WithTagsValue(context.Background(), "c", "a")
	ctx = ae.WithTagsValue(ctx, "b", "a")

	got := ae.TagsFromContext(ctx)
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("TagsFromContext (deduped+sorted) = %v, want %v", got, want)
	}
}

func TestBuilder_ContextPullsTagsIntoError(t *testing.T) {
	t.Parallel()

	ctx := ae.WithTagsValue(context.Background(), "ctx-tag")
	err := ae.NewC(ctx).Msg("x")

	got := ae.Tags(err)
	if !slices.Contains(got, "ctx-tag") {
		t.Errorf("Tags after NewC = %v, want to contain %q", got, "ctx-tag")
	}
}
