package ae

import (
	"context"
	"slices"
)

// ErrorTags defines an interface for errors that can provide a list of tags.
type ErrorTags interface {
	// ErrorTags returns a list of tags associated with the error.
	// Returns nil if no tags are set.
	ErrorTags() []string
}

// Tags extracts the list of tags from an error.
// If the error implements ErrorTags, returns its ErrorTags().
// Returns nil if err is nil or if the error does not implement ErrorTags.
func Tags(err error) []string {
	if err == nil {
		return nil
	}

	if ae, ok := err.(ErrorTags); ok {
		return ae.ErrorTags()
	}

	return nil
}

type tagKey struct{}

// WithTagsValue returns a new context with the given tags added to it.
// If the context already contains tags, the new tags are appended to the existing tags and de-duplicated.
func WithTagsValue(ctx context.Context, tags ...string) context.Context {
	existingTags, ok := ctx.Value(tagKey{}).([]string)
	if !ok {
		existingTags = []string{}
	}

	joinedTags := slices.Compact(
		slices.Sorted(
			slices.Values(append(existingTags, tags...)),
		),
	)

	return context.WithValue(ctx, tagKey{}, joinedTags)
}

func TagsFromContext(ctx context.Context) []string {
	tags, ok := ctx.Value(tagKey{}).([]string)
	if !ok {
		return nil
	}

	return tags
}
