package ae

import (
	"context"
	"maps"
)

// ErrorAttributes defines an interface for errors that can provide a map of attributes.
type ErrorAttributes interface {
	// ErrorAttributes returns a map of attributes associated with the error.
	// Returns an empty map non-nil if no attributes are set.
	ErrorAttributes() map[string]any
}

// Attributes extracts the map of attributes from an error.
// If the error implements ErrorAttributes, returns its Attributes().
// Returns an empty map if err is nil or if the error does not implement ErrorAttributes.
func Attributes(err error) map[string]any {
	if err == nil {
		return make(map[string]any)
	}

	if ae, ok := err.(ErrorAttributes); ok {
		attrs := ae.ErrorAttributes()
		if attrs != nil {
			return attrs
		}
	}

	return make(map[string]any)
}

type attributesKey struct{}

// WithAttribute creates a new context with the given attribute added to it.
// If the context already contains attributes, the new attribute is added to the existing attributes. Attributes
// with duplicate keys are overwritten.
func WithAttribute(ctx context.Context, key string, value any) context.Context {
	return WithAttributes(ctx, map[string]any{key: value})
}

// WithAttributes creates a new context with the given attributes added to it.
// If the context already contains attributes, the new attributes are added to the existing attributes. Attributes
// with duplicate keys are overwritten.
func WithAttributes(ctx context.Context, attrs map[string]any) context.Context {
	existingAttrs, ok := ctx.Value(tagKey{}).(map[string]any)
	if !ok {
		existingAttrs = make(map[string]any)
	}

	maps.Copy(existingAttrs, attrs)
	return context.WithValue(ctx, attributesKey{}, existingAttrs)
}

// AttributesFromContext extracts the attributes map from the given context.
// If the context contains no attributes, it returns an empty map.
func AttributesFromContext(ctx context.Context) map[string]any {
	attrs, ok := ctx.Value(tagKey{}).(map[string]any)
	if !ok {
		return make(map[string]any)
	}

	return attrs
}
