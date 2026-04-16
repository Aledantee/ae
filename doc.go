// Package ae is an opinionated error-handling library for Go.
//
// An [Ae] error carries more than a message: an optional user-facing copy
// and resolution hint, a timestamp, a string error code, a process exit
// code, OpenTelemetry trace and span IDs, a set of tags, a map of
// attributes, underlying causes, related (non-cause) errors, and captured
// goroutine stack traces. All of this metadata is accessible through
// small, single-purpose interfaces (see [ErrorMessage], [ErrorCode], and
// peers), so arbitrary error types that implement the relevant interface
// cooperate with the package-level accessors.
//
// # Getting started
//
// Three entry points cover most uses:
//
//   - [Msg] / [Msgf] — create a bare error with just a message.
//   - [Wrap] / [Wrapf] — add a new message on top of an existing error.
//   - [New] — open a [Builder] for anything richer (codes, tags, causes,
//     attributes, stacks, OpenTelemetry IDs).
//
// Each of the above has a context-aware sibling ([MsgC], [WrapC], [NewC],
// [MsgCf], [WrapCf]) that inherits attributes, tags, and OpenTelemetry
// trace/span IDs from the [context.Context] via [Builder.Context].
//
// # Layering
//
// The package is intentionally layered:
//
//   - The [Builder] is the source of truth. Its terminal methods
//     ([Builder.Msg], [Builder.Msgf], [Builder.UserMsg]) return an error —
//     note the signature change, you cannot chain further after them.
//   - The top-level convenience functions are thin wrappers around the
//     builder for the common cases.
//   - The [Printer] (via [Print] and [NewPrinter]) handles formatting.
//     Toggle individual fields with the [PrinterOption] family, or reach
//     for the [PrintVerbose] and [PrintCompact] presets.
//   - Accessor functions ([Message], [Code], [Tags], [Causes], …) read
//     metadata back out of any error — not just [Ae] values — via their
//     respective ErrorXxx interfaces.
//
// The errors sub-package (go.aledante.io/ae/errors) provides drop-in
// replacements for the standard library's [errors.New], [errors.Join],
// [errors.Is], [errors.As], and [errors.Unwrap] that produce ae errors.
//
// # Examples
//
// Runnable examples covering each major capability live under examples/
// in the repository. Run them with, for example, "go run ./examples/simple"
// or "go run ./examples/print".
package ae
