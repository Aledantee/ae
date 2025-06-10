package ae

// PrinterOption defines a function type that configures a Printer.
// It is used to customize the behavior of a Printer instance through functional options.
type PrinterOption func(p *Printer)

// WithUserMessage returns a PrinterOption that enables inclusion of user-friendly messages in the output.
func WithUserMessage() PrinterOption {
	return func(p *Printer) {
		p.userMsg = true
	}
}

// WithoutUserMessage returns a PrinterOption that disables inclusion of user-friendly messages in the output.
func WithoutUserMessage() PrinterOption {
	return func(p *Printer) {
		p.userMsg = false
	}
}

// WithHint returns a PrinterOption that enables inclusion of hint messages in the output.
// Hint messages provide suggestions for resolving the error.
func WithHint() PrinterOption {
	return func(p *Printer) {
		p.hint = true
	}
}

// WithoutHint returns a PrinterOption that disables inclusion of hint messages in the output.
func WithoutHint() PrinterOption {
	return func(p *Printer) {
		p.hint = false
	}
}

// WithTimestamp returns a PrinterOption that enables inclusion of error timestamps in the output.
func WithTimestamp() PrinterOption {
	return func(p *Printer) {
		p.timestamp = true
	}
}

// WithCode returns a PrinterOption that enables inclusion of error codes in the output.
func WithCode() PrinterOption {
	return func(p *Printer) {
		p.code = true
	}
}

// WithoutCode returns a PrinterOption that disables inclusion of error codes in the output.
func WithoutCode() PrinterOption {
	return func(p *Printer) {
		p.code = false
	}
}

// WithExitCode returns a PrinterOption that enables inclusion of exit codes in the output.
func WithExitCode() PrinterOption {
	return func(p *Printer) {
		p.exitCode = true
	}
}

// WithoutExitCode returns a PrinterOption that disables inclusion of exit codes in the output.
func WithoutExitCode() PrinterOption {
	return func(p *Printer) {
		p.exitCode = false
	}
}

// WithoutTimestamp returns a PrinterOption that disables inclusion of error timestamps in the output.
func WithoutTimestamp() PrinterOption {
	return func(p *Printer) {
		p.timestamp = false
	}
}

// WithStacks returns a PrinterOption that enables stack trace inclusion in the output.
func WithStacks() PrinterOption {
	return func(p *Printer) {
		p.stacks = true
	}
}

// WithoutStacks returns a PrinterOption that disables stack trace inclusion in the output.
func WithoutStacks() PrinterOption {
	return func(p *Printer) {
		p.stacks = false
	}
}

// WithJSON returns a PrinterOption that enables JSON formatting of the output.
func WithJSON() PrinterOption {
	return func(p *Printer) {
		p.json = true
	}
}

// WithoutJSON disables JSON formatting for the Printer, configuring it to produce plain text output instead.
func WithoutJSON() PrinterOption {
	return func(p *Printer) {
		p.json = false
	}
}

// WithIndent configures the Printer to use the specified number of spaces for indentation when formatting output.
// A minimum indentation of 1 is enforced.
func WithIndent(indent int) PrinterOption {
	if indent <= 0 {
		indent = 1
	}

	return func(p *Printer) {
		p.indent = indent
	}
}

// WithCauses returns a PrinterOption that enables inclusion of error causes in the output.
func WithCauses() PrinterOption {
	return func(p *Printer) {
		p.causes = true
	}
}

// WithRelated returns a PrinterOption that enables inclusion of related errors in the output.
func WithRelated() PrinterOption {
	return func(p *Printer) {
		p.related = true
	}
}

// WithoutRelated returns a PrinterOption that disables the inclusion of related errors in the printer's output.
func WithoutRelated() PrinterOption {
	return func(p *Printer) {
		p.related = false
	}
}

// WithoutCauses returns a PrinterOption that disables inclusion of error causes in the output.
func WithoutCauses() PrinterOption {
	return func(p *Printer) {
		p.causes = false
	}
}

// WithInfiniteDepth returns a PrinterOption that sets the error chain traversal depth to infinite.
// This means the printer will traverse the entire error chain regardless of depth.
func WithInfiniteDepth() PrinterOption {
	return func(p *Printer) {
		p.maxDepth = -1
	}
}

// WithMaxDepth returns a PrinterOption that sets a specific maximum depth for error chain traversal.
// The printer will stop traversing the error chain after reaching the specified depth.
func WithMaxDepth(depth int) PrinterOption {
	return func(p *Printer) {
		p.maxDepth = depth
	}
}

// WithColors returns a PrinterOption that enables colored output formatting.
func WithColors() PrinterOption {
	return func(p *Printer) {
		p.colors = true
	}
}

// WithoutColors returns a PrinterOption that disables colored output formatting.
func WithoutColors() PrinterOption {
	return func(p *Printer) {
		p.colors = false
	}
}

// WithTraceId returns a PrinterOption that enables inclusion of trace IDs in the output.
func WithTraceId() PrinterOption {
	return func(p *Printer) {
		p.traceId = true
	}
}

// WithoutTraceId returns a PrinterOption that disables inclusion of trace IDs in the output.
func WithoutTraceId() PrinterOption {
	return func(p *Printer) {
		p.traceId = false
	}
}

// WithSpanId returns a PrinterOption that enables inclusion of span IDs in the output.
func WithSpanId() PrinterOption {
	return func(p *Printer) {
		p.spanId = true
	}
}

// WithoutSpanId returns a PrinterOption that disables inclusion of span IDs in the output.
func WithoutSpanId() PrinterOption {
	return func(p *Printer) {
		p.spanId = false
	}
}

// WithTags returns a PrinterOption that enables inclusion of error tags in the output.
func WithTags() PrinterOption {
	return func(p *Printer) {
		p.tags = true
	}
}

// WithoutTags returns a PrinterOption that disables inclusion of error tags in the output.
func WithoutTags() PrinterOption {
	return func(p *Printer) {
		p.tags = false
	}
}

// WithAttributes returns a PrinterOption that enables inclusion of error attributes in the output.
func WithAttributes() PrinterOption {
	return func(p *Printer) {
		p.attributes = true
	}
}

// WithoutAttributes returns a PrinterOption that disables inclusion of error attributes in the output.
func WithoutAttributes() PrinterOption {
	return func(p *Printer) {
		p.attributes = false
	}
}

// WithVerbose returns a PrinterOption that enables all available output fields.
// This includes user messages, hints, timestamps, codes, exit codes, colors,
// trace IDs, span IDs, tags, attributes, causes, related errors, and stack traces.
func WithVerbose() PrinterOption {
	return withChained(
		WithUserMessage(),
		WithHint(),
		WithTimestamp(),
		WithCode(),
		WithExitCode(),
		WithColors(),
		WithTraceId(),
		WithSpanId(),
		WithTags(),
		WithAttributes(),
		WithCauses(),
		WithRelated(),
		WithStacks(),
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
