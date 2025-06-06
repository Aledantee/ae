package ae

import "context"

// errorBuilderKey is a private type used as a context key for storing ErrorBuilder instances.
type errorBuilderKey struct{}

// WithError stores an ErrorBuilder in the context.
// If the builder is nil, the original context is returned unchanged.
// This allows propagating error context through the call stack.
func WithError(ctx context.Context, builder *ErrorBuilder) context.Context {
	if builder != nil {
		return context.WithValue(ctx, errorBuilderKey{}, builder)
	}

	return ctx
}

// FromContext retrieves an ErrorBuilder from the context.
// If an ErrorBuilder exists in the context, it is returned.
// Otherwise, a new ErrorBuilder is created with the provided message
// and initialized with context values using the Context() method.
// Panics if the message is empty.
func FromContext(ctx context.Context, msg string) *ErrorBuilder {
	v := ctx.Value(errorBuilderKey{})
	if eb, ok := v.(*ErrorBuilder); ok && eb != nil {
		return eb
	}

	return New(msg).Context(ctx)
}
