# ae — Aledantee's error-handling library

An opinionated error library for Go. Every error carries structured
metadata (codes, tags, attributes, causes, related errors, stack traces,
OpenTelemetry trace/span IDs, hints, user-facing copy) and renders
through a dedicated text/JSON printer.

## Features

- **Rich metadata** — timestamps, error codes, exit codes, OTel trace/span IDs, tags, attributes, causes, related errors, captured goroutine stacks.
- **User-facing copy** — separate internal message, end-user message, and resolution hint.
- **Fluent builder** — chainable API that terminates with one of `Msg`, `Msgf`, or `UserMsg`.
- **Context integration** — `NewC(ctx)` / `FromC(ctx, err)` inherit attributes, tags, and OTel IDs from `context.Context`.
- **Recoverability** — `Builder.Fatal()` flags errors as non-recoverable; `IsRecoverable(err)` walks the chain.
- **Print or JSON** — semantic colored text for humans (auto-off outside a TTY); structured JSON for machines.
- **slog integration** — `*ae.Ae` implements `slog.LogValuer` for structured logging.
- **Drop-in `errors` replacement** — `go.aledante.io/ae/errors` re-implements `New`, `Join`, `Is`, `As`, `Unwrap` on top of ae.

## Installation

```bash
go get go.aledante.io/ae
```

## Quick start

### Bare message

```go
package main

import (
    "fmt"
    "go.aledante.io/ae"
)

func main() {
    err := ae.Msg("failed to process request")
    fmt.Println(err)
    // Output: failed to process request
}
```

`ae.Msg` / `ae.Msgf` / `ae.Wrap` / `ae.Wrapf` are the shortcut entry points; use the builder for anything richer.

### Wrapping a cause

```go
err := ae.New().
    Cause(errors.New("invalid input")).
    Tag("validation").
    Msg("failed to process request")

fmt.Println(err)
// Output: failed to process request: invalid input
```

`Error()` renders only the message chain. Tags, codes, hints, and other metadata surface through `ae.Print` or the extractor helpers.

### Stack capture

```go
func processData() error {
    return ae.New().
        Stack().                 // capture all goroutine stacks
        Code("DATA_FAILED").
        Msg("data processing failed")
}

ae.Print(processData())          // text (default)
ae.Print(processData(), ae.PrintJSON())
```

### Context integration

```go
func processWithContext(ctx context.Context, data string) error {
    // Attach attributes / tags / OTel IDs upstream…
    ctx = ae.WithAttribute(ctx, "data_length", len(data))

    // …and build an error that inherits them.
    return ae.NewC(ctx).
        Tag("processing").
        Msg("processing failed")
}
```

Use `ae.WrapC(ctx, msg, err)` / `ae.FromC(ctx, err)` when wrapping an existing error.

## Core concepts

### Builder

Everything non-terminal returns `Builder` for chaining. Terminal
methods — `Msg`, `Msgf`, `UserMsg` — return `error`, so nothing can be
chained after them.

```go
err := ae.New().
    Hint("Check your network connection").
    Code("NETWORK_ERROR").
    ExitCode(5).
    Tag("network").Tag("retryable").
    Attr("retry_count", 3).
    Attr("timeout", "30s").
    Cause(underlyingErr).
    Related(sideEffectErr).
    Stack().
    Now().
    UserMsg(
        "processing failed",                            // internal msg
        "Something went wrong. Please try again.",      // end-user msg
    )
```

Build-a-non-recoverable error: `ae.New().Fatal().Msg(...)` (shortcut for `.Recoverable(false)`).

### Extractors

Read metadata back out of **any** error. Each extractor honours its
respective `ErrorXxx` interface, so custom error types participate.

```go
ae.Message(err)       // ErrorMessage
ae.UserMessage(err)   // ErrorUserMessage
ae.Hint(err)          // ErrorHint
ae.Code(err)          // ErrorCode
ae.ExitCode(err)      // ErrorExitCode (recursive max over causes)
ae.Timestamp(err)     // ErrorTimestamp
ae.TraceId(err)       // ErrorTraceId
ae.SpanId(err)        // ErrorSpanId
ae.Tags(err)          // ErrorTags
ae.Attributes(err)    // ErrorAttributes
ae.Causes(err)        // ErrorCauses / Unwrap() []error / Unwrap() error / Cause() error
ae.Related(err)       // ErrorRelated
ae.Stacks(err)        // ErrorStacks
ae.IsRecoverable(err) // ErrorRecoverable (recursive, default true)
```

### Printing

```go
ae.Print(err)                              // text, default options
ae.Print(err, ae.PrintJSON())              // JSON, indented

// Presets
ae.Print(err, ae.PrintVerbose())           // every field on
ae.Print(err, ae.PrintCompact())           // terse high-signal set
```

Layout of the default text output (colors applied when stdout is a TTY):

```
[ERROR] {NETWORK_ERROR/5} processing failed [network, retryable]
  hint       Check your network connection
  shown      Something went wrong. Please try again.
  attrs      retry_count  3
             timeout      30s
  caused by  tcp reset by peer
```

Every printer option is toggled through the `Print*` / `NoPrint*`
family:

| Option | Default | Effect |
|---|---|---|
| `PrintJSON` / `NoPrintJSON` | text | Switch output format. |
| `PrintColors` / `NoPrintColors` | auto (TTY) | Force colors on/off. |
| `PrintIndent(n)` | 2 | Spaces per indent level. |
| `PrintDepth(n)` / `PrintDepthInfinite` | infinite | Cause-chain traversal depth. |
| `PrintUserMessage` / `NoPrintUserMessage` | verbose | Include the `shown` row when distinct from msg. |
| `PrintHint` / `NoPrintHint` | verbose | Include the `hint` row. |
| `PrintTimestamp` / `NoPrintTimestamp` | verbose | Include the `time` row. |
| `PrintCode` / `NoPrintCode` | verbose | Render `{CODE}` in the header. |
| `PrintExitCode` / `NoPrintExitCode` | verbose | Render `/N` (hidden when the default 1). |
| `PrintTags` / `NoPrintTags` | verbose | Include `[tag, tag]` in the header. |
| `PrintAttributes` / `NoPrintAttributes` | verbose | Include the `attrs` block. |
| `PrintCauses` / `NoPrintCauses` | verbose | Include the `caused by` block. |
| `PrintRelated` / `NoPrintRelated` | verbose | Include the `related` block. |
| `PrintStacks` / `NoPrintStacks` | verbose | Include the `stack` block. |
| `PrintTraceId` / `PrintSpanId` / `PrintOtel` | verbose | OTel IDs (PrintOtel = both). |
| `PrintFrameFilters(fn, …)` | ae+runtime hidden | Drop matching stack frames. |
| `PrintVerbose` / `PrintCompact` | verbose | Presets. |

### Distributed tracing

`Builder.Context(ctx)` — called by `NewC` / `FromC` — automatically
extracts OTel trace and span IDs from a span attached to the context:

```go
func processRequest(ctx context.Context) error {
    return ae.NewC(ctx).          // trace + span IDs pulled from ctx
        Msg("request processing failed")
}
```

Explicit IDs can be attached with `.TraceId(...)` / `.SpanId(...)`; call
them **after** `.Context(ctx)` if you need to override.

### Structured logging (slog)

`*ae.Ae` implements `slog.LogValuer`, so errors log as a structured
group when passed through `log/slog`:

```go
slog.Error("request failed", slog.Any("err", err))
```

Every populated field (msg, user_msg, hint, timestamp, code, exit_code,
tags, attributes, causes, related) surfaces as a slog attribute;
`recoverable` and `msg` are always present. Nested causes and related
errors render as their own sub-groups.

### errors sub-package

```go
import aeerrors "go.aledante.io/ae/errors"

aeerrors.New("…")                    // returns an ae.Ae error
aeerrors.Join(err1, nil, err2)       // nil-filtered, single passthrough, else bracketed
aeerrors.Is(err, target)             // proxies stdlib errors.Is
aeerrors.As(err, &target)            // proxies stdlib errors.As
aeerrors.Unwrap(err)                 // proxies stdlib errors.Unwrap
```

## Custom error types

Implement the relevant `ErrorXxx` interface — method names are prefixed
with `Error` — and the extractor helpers will pick up values:

```go
type myErr struct {
    msg  string
    code string
}

func (e myErr) Error() string      { return e.msg }
func (e myErr) ErrorMessage() string { return e.msg }
func (e myErr) ErrorCode() string    { return e.code }

code := ae.Code(myErr{code: "CUSTOM_001"})   // -> "CUSTOM_001"
```

## Examples

Runnable examples live under `examples/`. A good starting set:

- `examples/simple` — the minimum viable `Wrap` + `Print`.
- `examples/print` — ~12 formatted scenarios covering every printer path.
- `examples/chaining` — builder chaining with tags and causes.
- `examples/stack` — stack capture.
- `examples/context` — context-derived attributes.
- `examples/extract` — every extractor helper, side-by-side.
- `examples/errors` — drop-in `errors` package semantics.
- `examples/slog` — slog handler integration.
- `examples/exit` — using `ae.Exit` / `ae.PrintExit` with an exit code.

Run any of them with, for example, `go run ./examples/print`.
