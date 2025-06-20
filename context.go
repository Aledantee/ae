package ae

import "context"

// errorBuilderKey is a private type used as a context key for storing ErrorBuilder instances.
type errorBuilderKey struct{}

// WithError stores an ErrorBuilder in the context.
// If the builder is nil, the original context is returned unchanged.
// This allows propagating error context through the call stack.
func WithError(ctx context.Context, builder Builder) context.Context {
	return context.WithValue(ctx, errorBuilderKey{}, builder)
}

// FromContext retrieves an ErrorBuilder from the context.
// If an ErrorBuilder exists in the context, it is returned.
// Otherwise, a new ErrorBuilder is created and initialized
// with context values using the Context() method.
func FromContext(ctx context.Context) Builder {
	v := ctx.Value(errorBuilderKey{})
	if eb, ok := v.(Builder); ok {
		return eb
	}

	return New().Context(ctx)
}
