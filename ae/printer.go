package ae

import (
	"fmt"
)

// Printer provides functionality for formatting and printing errors with various options.
// It supports both plain text and JSON output formats, and can include stack traces
// and error causes in the output.
type Printer struct {
	// colors determines whether colored output is enabled.
	colors bool
	// json determines whether the output should be formatted as JSON
	json bool
	// indent is the number of spaces to indent by.
	indent int
	// maxDepth controls how deep to traverse the error chain when printing causes.
	// A negative value indicates infinite depth.
	maxDepth int

	// flags for error fields
	userMsg    bool
	hint       bool
	timestamp  bool
	code       bool
	exitCode   bool
	traceId    bool
	spanId     bool
	panId      bool
	tags       bool
	attributes bool
	causes     bool
	related    bool
	stacks     bool
}

// NewPrinter creates a new Printer with the given options.
// By default, the printer will:
//   - Include stack traces (stack = true)
//   - Output in plain text format (json = false)
//   - Include error causes (causes = true)
//   - Traverse the error chain infinitely (maxDepth = -1)
//
// These defaults can be overridden using PrinterOption functions.
func NewPrinter(opts ...PrinterOption) *Printer {
	opts = append([]PrinterOption{
		WithColors(),
		WithoutJSON(),
		WithIndent(2),
		WithVerbose(), // expands to all fields
	}, opts...)

	p := &Printer{}
	for _, opt := range append(opts, WithCauses()) {
		opt(p)
	}

	return p
}

func (p *Printer) PrettyPrint(err error) {
	p.Print(err)
}

// Print writes the formatted error to standard output.
// It uses the configured options (JSON format, stack traces, causes) to determine
// how to format the error.
func (p *Printer) Print(err error) {
	fmt.Println(p.Prints(err))
}

// Prints returns a string representation of the error based on the printer's configuration.
// If JSON output is enabled, it returns a JSON-formatted string.
// Otherwise, it returns a plain text representation.
// The output may include stack traces and error causes depending on the printer's configuration.
func (p *Printer) Prints(err error) string {
	if p.json {
		return p.printsJson(err, 0)
	} else {
		return p.printsText(err, 0)
	}
}
