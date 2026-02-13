package main

import (
	"log/slog"
	"os"

	"go.aledante.io/ae"
)

func main() {
	err := ae.New().
		Now().
		Attr("attr_0", "attr_0_value").
		Cause(
			ae.New().
				Tag("cause").
				Msg("cause"),
		).
		Msg("an error")

	textLogger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	jsonLogger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	textLogger.Error("oh no!", "error", err)
	jsonLogger.Error("oh no!", "error", err)
}
