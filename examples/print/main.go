// Example print exercises ae.Print and the printer-option surface across a
// dozen small scenarios: bare messages, hints, user-facing copy, timestamps,
// related errors, attributes, codes, exit codes, deep chains, trees, the
// PrintCompact preset, stack traces, no-color output, and JSON output.
// Each scenario is labeled so running the example shows each option's effect
// side-by-side.
package main

import (
	"errors"
	"fmt"
	"time"

	"go.aledante.io/ae"
)

func banner(s string) {
	fmt.Println()
	fmt.Println("===== " + s + " =====")
}

func main() {
	banner("1. bare message")
	ae.Print(ae.Msg("just a plain error"))

	banner("2. message + hint + user-facing copy")
	ae.Print(ae.New().
		Hint("try restarting the service").
		UserMsg(
			"auth token expired during refresh",
			"Your session expired. Please sign in again.",
		))

	banner("3. timestamp")
	ae.Print(ae.New().
		Timestamp(time.Date(2026, 4, 16, 9, 31, 4, 0, time.UTC)).
		Msg("timestamped error"))

	banner("4. related errors (non-cause)")
	ae.Print(ae.New().
		Related(
			errors.New("metrics flush failed"),
			errors.New("log write failed"),
		).
		Msg("cleanup had side effects"))

	banner("5. attributes, no colors (value must not be %v)")
	ae.Print(ae.New().
		Attr("user_id", 42).
		Attr("path", "/api/v1/login").
		Attr("retry_count", 3).
		Msg("attribute rendering"),
		ae.NoPrintColors())

	banner("6. code-only / exit-only / code+exit")
	ae.Print(ae.New().Code("E_AUTH").Msg("code only"))
	ae.Print(ae.New().ExitCode(2).Msg("exit only"))
	ae.Print(ae.New().Code("E_AUTH").ExitCode(77).Msg("code + exit"))

	banner("7. deep chain (5 levels)")
	chain := ae.Msg("level 0")
	for i := 1; i <= 4; i++ {
		chain = ae.Wrap(fmt.Sprintf("level %d", i), chain)
	}
	ae.Print(chain)

	banner("8. tree: multiple causes with nested children")
	ae.Print(ae.New().
		Cause(
			ae.New().
				Cause(
					ae.New().Tag("timeout").
						Cause(ae.New().Code("DEEP").Msg("deep nested error")).
						Msg("timeout"),
					ae.New().Tag("timeout").Msg("timeout2"),
				).
				Msg("database connection failed"),
			ae.New().Msg("cache miss"),
		).
		Code("AUTH_FAILED").
		ExitCode(77).
		TraceId("trace-abc123").
		SpanId("span-def456").
		Attr("path", "/api/v1/login").
		Attr("user_id", 42).
		Tag("security").
		Tag("auth").
		Hint("try rotating the service token").
		UserMsg(
			"authentication failed",
			"Something went wrong on our end. Please retry.",
		))

	banner("9. compact preset (no stacks, no timestamps, no trace)")
	ae.Print(ae.New().
		Code("E_NET").
		ExitCode(5).
		Tag("net").
		Attr("host", "api.example.com").
		Hint("retry with backoff").
		Cause(errors.New("tcp reset by peer")).
		Msg("network error"),
		ae.PrintCompact())

	banner("10. stack trace + attrs + cause")
	ae.Print(a())

	banner("11. no-colors variant of the kitchen-sink error")
	ae.Print(ae.New().
		Code("AUTH_FAILED").ExitCode(77).
		TraceId("trace-abc123").SpanId("span-def456").
		Tag("security").
		Attr("path", "/api/v1/login").
		Cause(errors.New("tcp reset by peer")).
		Hint("try rotating the service token").
		UserMsg("auth failed", "Please sign in again."),
		ae.NoPrintColors())

	banner("12. JSON output (machine-readable)")
	ae.Print(ae.New().
		Code("E_NET").
		ExitCode(5).
		Tag("net").
		Attr("host", "api.example.com").
		Cause(errors.New("tcp reset by peer")).
		Msg("network error"),
		ae.NoPrintColors(), ae.PrintJSON())
}

func a() error {
	return func() error {
		return b()
	}()
}

func b() error {
	return ae.New().
		Stack().
		Attr("operation", "insert").
		Attr("table", "users").
		Cause(errors.New("connection refused")).
		Code("DB_INSERT").
		Hint("check the database connection").
		Msg("failed to insert record")
}
