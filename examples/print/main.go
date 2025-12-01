package main

import "go.aledante.io/ae"

func main() {
	// Create simple error message
	err := ae.New().
		Cause(
			ae.New().
				Cause(
					ae.New().Tag("timeout").
						Cause(
							ae.New().Code("DEEP_NESTED_ERROR_CODE").Msg("deep nested error"),
						).
						Msg("timeout"),
					ae.New().Tag("timeout").Msg("timeout2"),
				).
				Msg("database connection failed"),
			ae.New().Msg("cause 2"),
		).
		Code("TEST_ERROR_CODE").
		TraceId("TRACE_ID").
		SpanId("SPAN_ID").
		Attr("key", "value").
		ExitCode(22).
		Tag("database").
		Hint("try restarting the application").
		UserMsg("something went wrong", "Oops! An error occurred.")

	// Print the error with default printer
	ae.Print(err)

	// Create printer with custom options using json formatting
	printer := ae.NewPrinter(
		ae.NoPrintColors(),
		ae.PrintJSON(),
	)

	// Print using custom printer
	//printer.Print(err)
	_ = printer
}
