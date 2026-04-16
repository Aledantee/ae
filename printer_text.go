package ae

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Color roles for text-mode rendering. Each call through Printer.fmt becomes a
// no-op when Printer.colors is false — the formatted string is returned verbatim.
// EnableColor is called on every instance so fatih/color does not second-guess
// our decision based on its own TTY detection: the Printer.colors flag is the
// single source of truth.
var (
	colBadge    = forceColor(color.New(color.FgRed, color.Bold))
	colMsg      = forceColor(color.New(color.FgRed, color.Bold))
	colCode     = forceColor(color.New(color.FgHiYellow))
	colBrace    = forceColor(color.New(color.FgYellow))
	colTag      = forceColor(color.New(color.FgHiMagenta))
	colBracket  = forceColor(color.New(color.FgMagenta))
	colLabel    = forceColor(color.New(color.FgCyan))
	colHint     = forceColor(color.New(color.FgHiCyan))
	colShown    = forceColor(color.New(color.FgWhite, color.Bold))
	colDim      = forceColor(color.New(color.FgHiBlack))
	colAttrKey  = forceColor(color.New(color.FgHiBlue))
	colAttrVal  = forceColor(color.New(color.FgHiGreen))
	colStackFn  = forceColor(color.New(color.FgHiYellow))
	colStackLoc = forceColor(color.New(color.FgHiBlack))
	colStackLn  = forceColor(color.New(color.FgYellow))
)

// forceColor returns c after calling EnableColor so fatih/color will emit ANSI
// regardless of the package-level NoColor/TTY detection. The Printer.colors
// flag still gates whether these instances get called at all.
func forceColor(c *color.Color) *color.Color {
	c.EnableColor()
	return c
}

const (
	// textLead is the indent before a section label.
	textLead = "  "
	// textLabelWidth is the padded width of the label column (fits "caused by").
	textLabelWidth = 9
	// textLabelGap is the spacing between label and value.
	textLabelGap = "  "
)

// textContinuationPrefix is the column where multi-line content and wrapped
// values begin — textLead + label column + gap, all spaces.
var textContinuationPrefix = textLead +
	strings.Repeat(" ", textLabelWidth) +
	textLabelGap

// fmt formats s with a...and colorizes with c when colors are enabled. It is the
// single funnel for every piece of text so color-on and color-off produce the
// same string content, only the ANSI wrapping changes.
func (p *Printer) fmt(format string, c *color.Color, a ...any) string {
	s := fmt.Sprintf(format, a...)
	if p.colors {
		return c.Sprint(s)
	}
	return s
}

// PrintErrorText renders err as a human-readable, labeled block. depth == 0
// marks the top-level call; nested calls (via causes / related) render through
// writeErrorTree using the inline form.
// The returned string is NOT newline-terminated.
func (p *Printer) PrintErrorText(err error, depth int) string {
	var sb strings.Builder
	p.writeHeader(&sb, err, depth == 0)
	p.writeSections(&sb, err, depth)
	return sb.String()
}

// writeHeader renders the first line: optional "[ERROR]" badge + inline summary.
func (p *Printer) writeHeader(sb *strings.Builder, err error, topLevel bool) {
	if topLevel {
		sb.WriteString(p.fmt("[ERROR]", colBadge))
		sb.WriteString(" ")
	}
	sb.WriteString(p.formatInlineError(err))
}

// formatInlineError renders the compact one-line form of an error:
//
//	{CODE/EXIT} message [tags]
//
// Used for both the top-level header and nested errors inside trees.
func (p *Printer) formatInlineError(err error) string {
	var sb strings.Builder

	code := ""
	exit := 0
	if p.code {
		code = Code(err)
	}
	if p.exitCode {
		// ExitCode(err) defaults to 1 for any non-nil error; that conventional
		// "error exit" is noise, so only render when the caller explicitly set
		// a distinct non-default value.
		if e := ExitCode(err); e > 1 {
			exit = e
		}
	}
	if code != "" || exit > 0 {
		sb.WriteString(p.fmt("{", colBrace))
		switch {
		case code != "" && exit > 0:
			sb.WriteString(p.fmt("%s", colCode, code))
			sb.WriteString(p.fmt("/", colBrace))
			sb.WriteString(p.fmt("%d", colCode, exit))
		case code != "":
			sb.WriteString(p.fmt("%s", colCode, code))
		default:
			sb.WriteString(p.fmt("exit ", colBrace))
			sb.WriteString(p.fmt("%d", colCode, exit))
		}
		sb.WriteString(p.fmt("}", colBrace))
		sb.WriteString(" ")
	}

	if msg := Message(err); msg != "" {
		sb.WriteString(p.fmt("%s", colMsg, msg))
	} else {
		sb.WriteString(p.fmt("(no message)", colDim))
	}

	if p.tags {
		if tags := Tags(err); len(tags) > 0 {
			sort.Strings(tags)
			sb.WriteString(" ")
			sb.WriteString(p.fmt("[", colBracket))
			for i, tag := range tags {
				if i > 0 {
					sb.WriteString(p.fmt(", ", colBracket))
				}
				sb.WriteString(p.fmt("%s", colTag, tag))
			}
			sb.WriteString(p.fmt("]", colBracket))
		}
	}

	return sb.String()
}

// writeSections emits the labeled rows below the header.
func (p *Printer) writeSections(sb *strings.Builder, err error, depth int) {
	if p.hint {
		if h := Hint(err); h != "" {
			p.writeRow(sb, "hint", p.fmt("%s", colHint, h))
		}
	}

	if p.userMsg {
		if u := UserMessage(err); u != "" && u != Message(err) {
			p.writeRow(sb, "shown", p.fmt("%s", colShown, u))
		}
	}

	if p.timestamp {
		if t := Timestamp(err); !t.IsZero() {
			p.writeRow(sb, "time", p.fmt("%s", colDim, t.Format(time.RFC3339)))
		}
	}

	if p.traceId || p.spanId {
		var parts []string
		if p.traceId {
			if id := TraceId(err); id != "" {
				parts = append(parts, p.fmt("%s", colDim, id))
			}
		}
		if p.spanId {
			if id := SpanId(err); id != "" {
				parts = append(parts,
					p.fmt("span ", colLabel)+p.fmt("%s", colDim, id))
			}
		}
		if len(parts) > 0 {
			p.writeRow(sb, "trace", strings.Join(parts, "  "))
		}
	}

	if p.attributes {
		if attrs := Attributes(err); len(attrs) > 0 {
			p.writeAttrs(sb, attrs)
		}
	}

	if p.causes && (p.maxDepth < 0 || depth < p.maxDepth) {
		if causes := Causes(err); len(causes) > 0 {
			p.writeErrorTree(sb, "caused by", causes, depth+1)
		}
	}

	if p.related {
		if related := Related(err); len(related) > 0 {
			p.writeErrorTree(sb, "related", related, depth+1)
		}
	}

	if p.stacks {
		if stacks := Stacks(err); len(stacks) > 0 {
			p.writeStacks(sb, stacks)
		}
	}
}

// writeRow writes a single labeled row on its own line.
func (p *Printer) writeRow(sb *strings.Builder, label, value string) {
	sb.WriteString("\n")
	sb.WriteString(p.labelPrefix(label))
	sb.WriteString(value)
}

// labelPrefix returns the prefix for the first line of a labeled block:
// leading indent + colored left-padded label + label gap. Its visual width
// matches textContinuationPrefix so subsequent lines align cleanly under it.
func (p *Printer) labelPrefix(label string) string {
	return textLead + p.fmt("%-*s", colLabel, textLabelWidth, label) + textLabelGap
}

// writeAttrs writes attributes sorted by key. The first pair shares the line
// with the "attrs" label so the block stays visually connected; subsequent
// pairs align under the first at textContinuationPrefix.
func (p *Printer) writeAttrs(sb *strings.Builder, attrs map[string]any) {
	keys := make([]string, 0, len(attrs))
	for k := range attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	maxKey := 0
	for _, k := range keys {
		if len(k) > maxKey {
			maxKey = len(k)
		}
	}

	for i, k := range keys {
		sb.WriteString("\n")
		if i == 0 {
			sb.WriteString(p.labelPrefix("attrs"))
		} else {
			sb.WriteString(textContinuationPrefix)
		}
		sb.WriteString(p.fmt("%-*s", colAttrKey, maxKey, k))
		sb.WriteString("  ")
		sb.WriteString(p.fmt("%v", colAttrVal, attrs[k]))
	}
}

// writeErrorTree prints a tree of errors (used for "caused by" and "related").
// The first top-level branch shares the line with the label; subsequent
// siblings and all nested children sit at textContinuationPrefix with the
// accumulated branch-accumulator for correct tree alignment.
//
// Glyph rules:
//   - Single cause at the top level: no glyph at all — reads as a plain
//     inline continuation of the label, so the reader never sees a line
//     that looks like it extends up into whatever section came before.
//   - Single cause nested (a "cause of a cause"): "└─" — standard tree
//     glyph, since it's clearly inside a tree and connects to the parent's
//     implicit stem.
//   - First of multiple at the top level: "┬─" (no up-stroke) so it cannot
//     be misread as continuing from whatever sat on the line above.
//   - First of multiple nested: "├─" — its up-stroke correctly lands on the
//     parent's down-stem, so the tree stays connected.
//   - Middle: "├─", last: "└─".
func (p *Printer) writeErrorTree(sb *strings.Builder, label string, errs []error, depth int) {
	p.writeErrorTreeRec(sb, label, errs, depth, "", true)
}

func (p *Printer) writeErrorTreeRec(sb *strings.Builder, label string, errs []error, depth int, branchAccum string, topLevel bool) {
	single := len(errs) == 1

	for i, e := range errs {
		isFirst := i == 0
		isLast := i == len(errs)-1

		var glyph, nextAccum string
		switch {
		case single && topLevel:
			// Top-level single cause sits flush with the label column —
			// no glyph so nothing looks like it connects up into the
			// section above the "caused by" label.
			glyph = ""
			nextAccum = branchAccum + "   "
		case single:
			// Nested single cause — a "cause of a cause". Standard tree
			// convention: use └─.
			glyph = p.fmt("└─ ", colDim)
			nextAccum = branchAccum + "   "
		case isFirst && topLevel:
			// First of multiple at top level — T-down glyph has no up-stroke
			// so it never reads as continuing from the line above.
			glyph = p.fmt("┬─ ", colDim)
			nextAccum = branchAccum + p.fmt("│  ", colDim)
		case isLast:
			glyph = p.fmt("└─ ", colDim)
			nextAccum = branchAccum + "   "
		default:
			glyph = p.fmt("├─ ", colDim)
			nextAccum = branchAccum + p.fmt("│  ", colDim)
		}

		sb.WriteString("\n")
		if label != "" && isFirst {
			sb.WriteString(p.labelPrefix(label))
		} else {
			sb.WriteString(textContinuationPrefix)
		}
		sb.WriteString(branchAccum)
		sb.WriteString(glyph)
		sb.WriteString(p.formatInlineError(e))

		if p.hint {
			if h := Hint(e); h != "" {
				sb.WriteString(" ")
				sb.WriteString(p.fmt("(%s)", colHint, h))
			}
		}

		if p.maxDepth < 0 || depth < p.maxDepth {
			if nested := Causes(e); len(nested) > 0 {
				p.writeErrorTreeRec(sb, "", nested, depth+1, nextAccum, false)
			}
		}
	}
}

// writeStacks prints captured goroutine stacks. The first goroutine header
// shares the line with the "stack" label; frames indent two columns further
// so the hierarchy is visually obvious. Frames are filtered through
// p.frameFilters — any frame for which a filter returns true is dropped, and
// a goroutine whose frames are all filtered out is omitted entirely.
func (p *Printer) writeStacks(sb *strings.Builder, stacks []*Stack) {
	frameIndent := textContinuationPrefix + "  "

	first := true
	for _, st := range stacks {
		frames := p.filterFrames(st.Frames)
		if len(frames) == 0 {
			continue
		}

		sb.WriteString("\n")
		if first {
			sb.WriteString(p.labelPrefix("stack"))
			first = false
		} else {
			sb.WriteString(textContinuationPrefix)
		}
		sb.WriteString(p.fmt("goroutine %d (%s)", colDim, st.ID, st.State))
		if st.Locked {
			sb.WriteString(p.fmt(" [locked]", colDim))
		}
		if st.Wait > 0 {
			sb.WriteString(p.fmt(" [wait=%s]", colDim, st.Wait))
		}

		maxFn := 0
		for _, f := range frames {
			if len(f.Func) > maxFn {
				maxFn = len(f.Func)
			}
		}
		for _, f := range frames {
			sb.WriteString("\n")
			sb.WriteString(frameIndent)
			sb.WriteString(p.fmt("%-*s", colStackFn, maxFn, f.Func))
			sb.WriteString(p.fmt("  at  ", colDim))
			sb.WriteString(p.fmt("%s", colStackLoc, f.File))
			sb.WriteString(p.fmt(":", colDim))
			sb.WriteString(p.fmt("%d", colStackLn, f.Line))
		}

		if st.FramesElided {
			sb.WriteString("\n")
			sb.WriteString(frameIndent)
			sb.WriteString(p.fmt("(frames elided)", colDim))
		}
	}
}

// filterFrames returns the subset of frames that survive every predicate in
// p.frameFilters — a frame is kept only when every filter returns false.
func (p *Printer) filterFrames(frames []*StackFrame) []*StackFrame {
	if len(p.frameFilters) == 0 {
		return frames
	}
	kept := make([]*StackFrame, 0, len(frames))
	for _, f := range frames {
		drop := false
		for _, filter := range p.frameFilters {
			if filter(f) {
				drop = true
				break
			}
		}
		if !drop {
			kept = append(kept, f)
		}
	}
	return kept
}
