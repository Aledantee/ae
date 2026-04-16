package ae

import (
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
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
	tags       bool
	attributes bool
	causes     bool
	related    bool
	stacks     bool

	// frameFilters is a list of predicates. A stack frame is dropped from the
	// rendered output when any filter returns true. The default set hides
	// internal ae/runtime frames; callers extend the list via PrintFrameFilters.
	frameFilters []func(frame *StackFrame) bool
}

// NewPrinter creates a new Printer with the given options.
//
// Defaults:
//   - Colors enabled when stdout is a terminal, disabled otherwise (via fatih/color.NoColor).
//   - Plain text output (json = false).
//   - Verbose field set (PrintVerbose enables every field).
//   - Infinite error-chain traversal (maxDepth = -1).
//   - Indent = 2.
//
// Defaults can be overridden by passing options. Later options win over earlier ones,
// so user-supplied options always override the built-in defaults.
func NewPrinter(opts ...PrinterOption) *Printer {
	colorsDefault := PrintColors()
	if color.NoColor {
		colorsDefault = NoPrintColors()
	}

	opts = append([]PrinterOption{
		colorsDefault,
		NoPrintJSON(),
		PrintIndent(2),
		PrintVerbose(),
		PrintDepthInfinite(),
	}, opts...)

	p := &Printer{
		frameFilters: []func(frame *StackFrame) bool{
			hideInternalFrames,
		},
	}
	for _, opt := range opts {
		opt(p)
	}

	return p
}

// hideInternalFrames is the default frame filter applied by NewPrinter. It
// drops frames whose function names belong to this library or Go's runtime
// stack-capture helpers, keeping the printed trace focused on user code.
func hideInternalFrames(frame *StackFrame) bool {
	if frame == nil {
		return true
	}
	return strings.HasPrefix(frame.Func, "go.aledante.io/ae") ||
		strings.HasPrefix(frame.Func, "runtime/debug")
}

// Print is a shortcut for NewPrinter(opts...).Print(err).
func Print(err error, opts ...PrinterOption) {
	NewPrinter(opts...).Print(err)
}

// PrettyPrint is an alias for Print.
func (p *Printer) PrettyPrint(err error) {
	p.Print(err)
}

// Print writes the formatted error to standard output followed by a single newline.
func (p *Printer) Print(err error) {
	p.Fprint(os.Stdout, err)
}

// Fprint writes the formatted error to w followed by a single newline.
func (p *Printer) Fprint(w io.Writer, err error) {
	io.WriteString(w, p.Prints(err))
	io.WriteString(w, "\n")
}

// Prints returns a string representation of the error based on the printer's configuration.
// If JSON output is enabled, it returns a JSON-formatted string.
// Otherwise, it returns a plain text representation.
// The returned string is NOT newline-terminated.
func (p *Printer) Prints(err error) string {
	if p.json {
		return p.printsJson(err, 0)
	}
	return p.PrintErrorText(err, 0)
}
