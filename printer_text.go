package ae

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	colSymbol = color.New(color.FgWhite)
	colMsg    = color.New(color.FgRed).
			Add(color.Bold)
	colUsrMsg  = color.New(color.FgWhite)
	colTag     = color.New(color.FgHiMagenta)
	colAttrKey = color.New(color.FgHiBlue)
	colAttrVal = color.New(color.FgHiGreen)
	colHint    = color.New(color.FgHiCyan)
	colCode    = color.New(color.FgHiYellow)
)

// formatErrorLine formats a single error line with code, exit code, message, tags, and hint
func (p *Printer) formatErrorLine(err error) string {
	var sb strings.Builder

	// Print code if enabled and available
	codeOpen := false
	if p.code {
		if code := Code(err); code != "" {
			sb.WriteString(p.fmt("{", colCode))
			sb.WriteString(p.fmt(code, colCode))
			codeOpen = true
		}
	}

	// Print exit code if enabled
	if p.exitCode {
		if exitCode := ExitCode(err); exitCode > 0 {
			sb.WriteString(p.fmt(fmt.Sprintf("/%d", exitCode), colCode))
		}
	}

	if codeOpen {
		sb.WriteString(p.fmt("} ", colCode))
	}

	// Print message
	sb.WriteString(p.fmt(Message(err), colMsg))

	// Print tags if enabled and available
	if p.tags {
		if tags := Tags(err); len(tags) > 0 {
			sb.WriteString(" [")
			for i, tag := range tags {
				if i > 0 {
					sb.WriteString(", ")
				}
				sb.WriteString(p.fmt(tag, colTag))
			}
			sb.WriteString("]")
		}
	}

	// Print hint if enabled and available
	if p.hint {
		if hint := Hint(err); hint != "" {
			sb.WriteString(" (")
			sb.WriteString(p.fmt(hint, colHint))
			sb.WriteRune(')')
		}
	}

	return sb.String()
}

// formatAttributeLine formats a single attribute line
func (p *Printer) formatAttributeLine(indent int, key string, value any) string {
	return fmt.Sprintf("%s-> %s: %s", strings.Repeat(" ", indent), p.fmt(key, colAttrKey), p.fmt("%v", colAttrVal, value))
}

// printErrorCauses recursively prints the error causes with proper tree structure
func (p *Printer) printErrorCauses(causes []error, depth int, sb *strings.Builder, prefix string) {
	if len(causes) == 0 {
		return
	}

	for i, cause := range causes {
		isLast := i == len(causes)-1
		branch := "├─ "
		nextPrefix := prefix + "│  "
		if isLast {
			branch = "└─ "
			nextPrefix = prefix + "   "
		}
		sb.WriteString(prefix)
		sb.WriteString(branch)
		sb.WriteString(p.formatErrorLine(cause))
		sb.WriteString("\n")

		// Recursively print nested causes
		if nestedCauses := Causes(cause); len(nestedCauses) > 0 && (p.maxDepth < 0 || depth < p.maxDepth) {
			p.printErrorCauses(nestedCauses, depth+1, sb, nextPrefix)
		}
	}
}

// PrintErrorText is the main entry for printing an error and its details as text
func (p *Printer) PrintErrorText(err error, depth int) string {
	var sb strings.Builder
	sb.WriteString(p.formatErrorLine(err))

	attrs := make(map[string]any)
	if p.attributes {
		attrs = Attributes(err)
	}

	if p.traceId {
		if traceId := TraceId(err); traceId != "" {
			attrs["Trace ID"] = traceId
		}
	}
	if p.spanId {
		if spanId := SpanId(err); spanId != "" {
			attrs["Span ID"] = spanId
		}
	}

	if len(attrs) > 0 {
		for k, v := range attrs {
			sb.WriteRune('\n')
			sb.WriteString(p.formatAttributeLine(p.indent, k, v))
		}
	}

	// Print causes if enabled and available
	if p.causes && (p.maxDepth < 0 || depth < p.maxDepth) {
		if causes := Causes(err); len(causes) > 0 {
			if depth == 0 {
				sb.WriteString("\nCauses:\n")
			}
			p.printErrorCauses(causes, depth+1, &sb, "")
		}
	}

	// Print stack traces if enabled and available
	if p.stacks {
		if stacks := Stacks(err); len(stacks) > 0 {
			sb.WriteString(p.fmt("Stack Traces:\n", colCode))
			for i, stack := range stacks {
				prefix := "└─ "
				if i < len(stacks)-1 {
					prefix = "├─ "
				}
				sb.WriteString(strings.Repeat(" ", p.indent))
				sb.WriteString(prefix)
				sb.WriteString(p.fmt(fmt.Sprintf("Goroutine %d (%s):\n", stack.ID, stack.State), colCode))
				for j, frame := range stack.Frames {
					prefix := "└─ "
					if j < len(stack.Frames)-1 {
						prefix = "├─ "
					}
					sb.WriteString(strings.Repeat(" ", p.indent*2))
					sb.WriteString(prefix)
					sb.WriteString(p.fmt(fmt.Sprintf("%s\n", frame.Func), colCode))
					sb.WriteString(strings.Repeat(" ", p.indent*2))
					sb.WriteString("   ")
					sb.WriteString(p.fmt(fmt.Sprintf("at %s:%d\n", frame.File, frame.Line), colCode))
				}
			}
		}
	}

	return sb.String()
}

func (p *Printer) fmt(s string, c *color.Color, fmt ...any) string {
	if p.colors {
		return c.Sprintf(s, fmt...)
	}

	return s
}
