package main

import (
	"context"

	"github.com/aledantee/ae"
)

func main() {
	ae.Exit(a())
}

func a() error {
	ctx := ae.WithAttributeValue(context.Background(), "test", "value")
	return b(ctx)
}

func b(ctx context.Context) error {
	return ae.NewC(ctx).Msg("error with context")
}
