package ae

import "context"

type errorKey struct{}

func WithError(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, errorKey{}, err)
}

func FromContext(ctx context.Context) error {
	return ctx.Value(errorKey{}).(error)
}
