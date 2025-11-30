package ae

// PrinterOption defines a function type that configures a Printer.
// It is used to customize the behavior of a Printer instance through functional options.
type PrinterOption func(p *Printer)

// PrintUserMessage returns a PrinterOption that enables inclusion of user-friendly messages in the output.
func PrintUserMessage() PrinterOption {
	return func(p *Printer) {
		p.userMsg = true
	}
}

// NoPrintUserMessage returns a PrinterOption that disables inclusion of user-friendly messages in the output.
func NoPrintUserMessage() PrinterOption {
	return func(p *Printer) {
		p.userMsg = false
	}
}

// PrintHint returns a PrinterOption that enables inclusion of hint messages in the output.
// Hint messages provide suggestions for resolving the error.
func PrintHint() PrinterOption {
	return func(p *Printer) {
		p.hint = true
	}
}

// NoPrintHint returns a PrinterOption that disables inclusion of hint messages in the output.
func NoPrintHint() PrinterOption {
	return func(p *Printer) {
		p.hint = false
	}
}

// PrintTimestamp returns a PrinterOption that enables inclusion of error timestamps in the output.
func PrintTimestamp() PrinterOption {
	return func(p *Printer) {
		p.timestamp = true
	}
}

// PrintCode returns a PrinterOption that enables inclusion of error codes in the output.
func PrintCode() PrinterOption {
	return func(p *Printer) {
		p.code = true
	}
}

// NoPrintCode returns a PrinterOption that disables inclusion of error codes in the output.
func NoPrintCode() PrinterOption {
	return func(p *Printer) {
		p.code = false
	}
}

// PrintExitCode returns a PrinterOption that enables inclusion of exit codes in the output.
func PrintExitCode() PrinterOption {
	return func(p *Printer) {
		p.exitCode = true
	}
}

// NoPrintExitCode returns a PrinterOption that disables inclusion of exit codes in the output.
func NoPrintExitCode() PrinterOption {
	return func(p *Printer) {
		p.exitCode = false
	}
}

// NoPrintTimestamp returns a PrinterOption that disables inclusion of error timestamps in the output.
func NoPrintTimestamp() PrinterOption {
	return func(p *Printer) {
		p.timestamp = false
	}
}

// PrintStacks returns a PrinterOption that enables stack trace inclusion in the output.
func PrintStacks() PrinterOption {
	return func(p *Printer) {
		p.stacks = true
	}
}

// NoPrintStacks returns a PrinterOption that disables stack trace inclusion in the output.
func NoPrintStacks() PrinterOption {
	return func(p *Printer) {
		p.stacks = false
	}
}

// PrintJSON returns a PrinterOption that enables JSON formatting of the output.
func PrintJSON() PrinterOption {
	return func(p *Printer) {
		p.json = true
	}
}

// NoPrintJSON disables JSON formatting for the Printer, configuring it to produce plain text output instead.
func NoPrintJSON() PrinterOption {
	return func(p *Printer) {
		p.json = false
	}
}

// PrintIndent configures the Printer to use the specified number of spaces for indentation when formatting output.
// A minimum indentation of 1 is enforced.
func PrintIndent(indent int) PrinterOption {
	if indent <= 0 {
		indent = 1
	}

	return func(p *Printer) {
		p.indent = indent
	}
}

// PrintCauses returns a PrinterOption that enables inclusion of error causes in the output.
func PrintCauses() PrinterOption {
	return func(p *Printer) {
		p.causes = true
	}
}

// NoPrintCauses returns a PrinterOption that disables inclusion of error causes in the output.
func NoPrintCauses() PrinterOption {
	return func(p *Printer) {
		p.causes = false
	}
}

// PrintRelated returns a PrinterOption that enables inclusion of related errors in the output.
func PrintRelated() PrinterOption {
	return func(p *Printer) {
		p.related = true
	}
}

// NoPrintRelated returns a PrinterOption that disables the inclusion of related errors in the printer's output.
func NoPrintRelated() PrinterOption {
	return func(p *Printer) {
		p.related = false
	}
}

// PrintDepthInfinite returns a PrinterOption that sets the error chain traversal depth to infinite.
// This means the printer will traverse the entire error chain regardless of depth.
func PrintDepthInfinite() PrinterOption {
	return func(p *Printer) {
		p.maxDepth = -1
	}
}

// PrintDepth returns a PrinterOption that sets a specific maximum depth for error chain traversal.
// The printer will stop traversing the error chain after reaching the specified depth.
func PrintDepth(depth int) PrinterOption {
	return func(p *Printer) {
		p.maxDepth = depth
	}
}

// PrintColors returns a PrinterOption that enables colored output formatting.
func PrintColors() PrinterOption {
	return func(p *Printer) {
		p.colors = true
	}
}

// NoPrintColors returns a PrinterOption that disables colored output formatting.
func NoPrintColors() PrinterOption {
	return func(p *Printer) {
		p.colors = false
	}
}

// PrintOtel returns a PrinterOption that enables inclusion of OTEL information in the output.
func PrintOtel() PrinterOption {
	return func(p *Printer) {
		p.traceId = true
	}
}

// NoPrintOtel returns a PrinterOption that disables inclusion of OTEL information in the output.
func NoPrintOtel() PrinterOption {
	return func(p *Printer) {
		p.traceId = false
	}
}

// PrintTags returns a PrinterOption that enables inclusion of error tags in the output.
func PrintTags() PrinterOption {
	return func(p *Printer) {
		p.tags = true
	}
}

// NoPrintTags returns a PrinterOption that disables inclusion of error tags in the output.
func NoPrintTags() PrinterOption {
	return func(p *Printer) {
		p.tags = false
	}
}

// PrintAttributes returns a PrinterOption that enables inclusion of error attributes in the output.
func PrintAttributes() PrinterOption {
	return func(p *Printer) {
		p.attributes = true
	}
}

// NoPrintAttributes returns a PrinterOption that disables inclusion of error attributes in the output.
func NoPrintAttributes() PrinterOption {
	return func(p *Printer) {
		p.attributes = false
	}
}

// PrintVerbose returns a PrinterOption that enables all available output fields.
// This includes user messages, hints, timestamps, codes, exit codes, colors,
// trace IDs, span IDs, tags, attributes, causes, related errors, and stack traces.
func PrintVerbose() PrinterOption {
	return withChained(
		PrintHint(),
		PrintTimestamp(),
		PrintCode(),
		PrintExitCode(),
		PrintColors(),
		PrintOtel(),
		PrintTags(),
		PrintAttributes(),
		PrintCauses(),
		PrintRelated(),
		PrintStacks(),
	)
}

// PrintCompact returns a PrinterOption that enables a minimal set of commonly useful output fields.
func PrintCompact() PrinterOption {
	return withChained(
		PrintHint(),
		PrintCode(),
		PrintExitCode(),
		PrintAttributes(),
		PrintTags(),
		PrintCauses(),
		PrintRelated(),
	)
}

// withChained combines multiple PrinterOptions into a single option that applies all of them.
func withChained(opts ...PrinterOption) PrinterOption {
	return func(p *Printer) {
		for _, o := range opts {
			o(p)
		}
	}
}
