package ae_test

import (
	"time"

	"go.aledante.io/ae"
)

// stubErr is a configurable test double that selectively implements the various
// ErrorX interfaces exported by the ae package. Tests use it to verify that
// extractor functions and Builder.From(err) pick values up from any error that
// satisfies the interface, not just *ae.Ae.
type stubErr struct {
	msg       string
	userMsg   string
	code      string
	exitCode  int
	hint      string
	traceId   string
	spanId    string
	tags      []string
	attrs     map[string]any
	causes    []error
	related   []error
	stacks    []*ae.Stack
	timestamp time.Time
}

func (s stubErr) Error() string                  { return s.msg }
func (s stubErr) ErrorMessage() string           { return s.msg }
func (s stubErr) ErrorUserMessage() string       { return s.userMsg }
func (s stubErr) ErrorCode() string              { return s.code }
func (s stubErr) ErrorExitCode() int             { return s.exitCode }
func (s stubErr) ErrorHint() string              { return s.hint }
func (s stubErr) ErrorTraceId() string           { return s.traceId }
func (s stubErr) ErrorSpanId() string            { return s.spanId }
func (s stubErr) ErrorTags() []string            { return s.tags }
func (s stubErr) ErrorAttributes() map[string]any { return s.attrs }
func (s stubErr) ErrorCauses() []error           { return s.causes }
func (s stubErr) ErrorRelated() []error          { return s.related }
func (s stubErr) ErrorStacks() []*ae.Stack       { return s.stacks }
func (s stubErr) ErrorTimestamp() time.Time      { return s.timestamp }

// multiUnwrapErr exercises the `Unwrap() []error` branch of ae.Causes and of
// Builder.CauseUnwrap / Builder.RelatedUnwrap.
type multiUnwrapErr struct {
	msg  string
	errs []error
}

func (m multiUnwrapErr) Error() string   { return m.msg }
func (m multiUnwrapErr) Unwrap() []error { return m.errs }

// singleUnwrapErr exercises the `Unwrap() error` branch of ae.Causes.
type singleUnwrapErr struct {
	msg   string
	inner error
}

func (s singleUnwrapErr) Error() string { return s.msg }
func (s singleUnwrapErr) Unwrap() error { return s.inner }

// pkgErrorsStyleErr exercises the `Cause() error` branch of ae.Causes
// (github.com/pkg/errors-style errors).
type pkgErrorsStyleErr struct {
	msg   string
	inner error
}

func (p pkgErrorsStyleErr) Error() string { return p.msg }
func (p pkgErrorsStyleErr) Cause() error  { return p.inner }

// plainErr is a bare error implementing only the builtin error interface.
type plainErr struct{ msg string }

func (p plainErr) Error() string { return p.msg }
