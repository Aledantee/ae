package ae

import (
	"encoding/json"
	"strings"
)

type jsonError struct {
	Message     string         `json:"message,omitempty"`
	UserMessage string         `json:"user_message,omitempty"`
	Hint        string         `json:"hint,omitempty"`
	Code        string         `json:"code,omitempty"`
	ExitCode    int            `json:"exit_code,omitempty"`
	TraceId     string         `json:"trace_id,omitempty"`
	SpanId      string         `json:"span_id,omitempty"`
	Tags        []string       `json:"tags,omitempty"`
	Attrs       map[string]any `json:"attrs,omitempty"`
	Causes      []jsonError    `json:"causes,omitempty"`
	Related     []jsonError    `json:"related,omitempty"`
	Stacks      []*Stack       `json:"stacks,omitempty"`
}

func (p *Printer) printsJson(err error, depth int) string {
	jsonErr := p.toJsonError(err, depth)
	jsonStr, _ := json.MarshalIndent(jsonErr, "", strings.Repeat(" ", p.indent))

	return string(jsonStr)
}

func (p *Printer) toJsonError(err error, depth int) jsonError {
	var (
		causes  []jsonError
		related []jsonError
	)

	if p.maxDepth < 0 || depth < p.maxDepth {
		for _, c := range Causes(err) {
			causes = append(causes, p.toJsonError(c, depth+1))
		}
		for _, r := range Related(err) {
			related = append(related, p.toJsonError(r, depth+1))
		}
	}

	je := jsonError{
		Message:     Message(err),
		UserMessage: UserMessage(err),
		Hint:        Hint(err),
		Code:        Code(err),
		ExitCode:    ExitCode(err),
		TraceId:     TraceId(err),
		SpanId:      SpanId(err),
		Tags:        Tags(err),
		Attrs:       Attributes(err),
		Causes:      causes,
		Related:     related,
		Stacks:      Stacks(err),
	}

	return je
}
