# AE - **A**ledantee's **E**rror Handling Library

This is another error handling library for Go, but within my control.

## Features

- **Rich Error Metadata**: Timestamps, error codes, exit codes, distributed tracing IDs, tags, and custom attributes
- **User-Friendly Messages**: Separate internal and user-facing error messages with resolution hints
- **Fluent API**: Clean, chainable builder pattern for constructing errors
- **Context Integration**: Built-in support for propagating error context through call chains
- **Multiple Output Formats**: JSON and formatted text output with extensive customization options

## Installation

```bash
go get github.com/aledantee/ae
```

## Quick Start

### Basic Error Creation

```go
package main

import (
    "errors"
    "github.com/aledantee/ae"
)

func main() {
    // Create a simple error
    err := ae.New().
        Msg("failed to process request").
        Code("PROCESS_FAILED").
        ExitCode(1)
    
    fmt.Println(err)
    // Output: failed to process request
}
```

### Error Chaining

```go
// Chain errors with causes
err := ae.New().
    Cause(errors.New("database connection failed")).
    Tag("database").
    Tag("connection").
    Msg("failed to process request")

fmt.Println(err)
// Output: failed to process request: database connection failed {database, connection}
```

### Stack Trace Capture

```go
func processData() error {
    return ae.New().
        Stack().  // Capture current stack trace
        Msg("data processing failed")
}

func main() {
    err := processData()
    
    // Print with stack trace in JSON format
    ae.Print(err, ae.WithJSON())
}
```

### Context Integration

```go
func processWithContext(ctx context.Context, data string) error {
    // Add an error builder to the context
    ctx = ae.WithError(ctx, ae.New().Tag("processing"))
    
    // Later in the call chain...
    return ae.FromContext(ctx).
        Attr("data_length", len(data)).
        Msg("processing failed")
}
```

## Core Concepts

### Error Builder

The `ae.Builder` provides a fluent interface for constructing rich errors:

```go
err := ae.New().
    UserMsg("Something went wrong. Please try again.").
    Hint("Check your network connection").
    Code("NETWORK_ERROR").
    ExitCode(1).
    Tag("network").
    Tag("retryable").
    Attr("retry_count", 3).
    Attr("timeout", "30s").
    Cause(underlyingError).
    Related(relatedError1, relatedError2).
    Stack().
    Now().  // Set current timestamp
    Msg("internal error message") // Msg is the terminal operation converting the builder to an error
```

### Error Extraction

Extract specific information from any error:

```go
// Extract basic information
message := ae.Message(err) // ae.ErrorMessage interface
userMessage := ae.UserMessage(err) // ae.ErrorUserMessage interface
hint := ae.Hint(err) // ae.ErrorHint interface
code := ae.Code(err) // ae.ErrorCode interface
exitCode := ae.ExitCode(err) // ae.ErrorExitCode interface

// Extract tracing information
traceID := ae.TraceId(err) // ae.ErrorTraceId interface
spanID := ae.SpanId(err) // ae.ErrorSpanId interface

// Extract metadata
tags := ae.Tags(err) // ErrorTags interface
attributes := ae.Attributes(err) // ErrorAttributes interface

// Extract error relationships
causes := ae.Causes(err) // ErrorCauses interface
related := ae.Related(err) // ErrorRelated interface
stacks := ae.Stacks(err) // ErrorStacks interface
```

### Error Printing

Customize error output with extensive options:

```go
// Basic printing
ae.Print(err)

// JSON output
ae.Print(err, ae.WithJSON())

// Custom formatting
ae.Print(err, 
    ae.WithJSON(),
    ae.WithIndent(4),
    ae.WithoutColors(),
    ae.WithVerbose(),  // Include all fields
    ae.WithInfiniteDepth(),  // Traverse all error chains
    ae.WithStacks(),  // Include stack traces
    ae.WithCauses(),  // Include error causes
    ae.WithAttributes(),  // Include custom attributes
)
```

### Context Propagation

Propagate error context through your call stack:

```go
func middleware(ctx context.Context, next func(context.Context) error) error {
    // Add context to the request
    ctx = ae.WithError(ctx, ae.New().Tag("middleware"))
    
    err := next(ctx)
    if err != nil {
        // Enrich the error with context
        return ae.FromContext(ctx).
            Cause(err).
            Msg("middleware failed")
    }
    
    return nil
}
```

### Distributed Tracing

Integrate with OpenTelemetry for distributed tracing:

```go
import (
    "go.opentelemetry.io/otel/trace"
)

func processRequest(ctx context.Context) error {
    span := trace.SpanFromContext(ctx)
    
    return ae.New().
        TraceId(span.SpanContext().TraceID().String()).
        SpanId(span.SpanContext().SpanID().String()).
        Context(ctx). // Also extracts trace and span IDs from the context.
        Msg("request processing failed")
}
```

## Advanced Usage

### Error Enrichment

Add related errors to existing errors:

```go
originalErr := errors.New("original error")
relatedErr1 := errors.New("related error 1")
relatedErr2 := errors.New("related error 2")

enrichedErr := ae.AddRelated(originalErr, relatedErr1, relatedErr2)
```

### Custom Error Types

The library works with any error type through interface implementations:

```go
type CustomError struct {
    message string
    code    string
}

func (e CustomError) Error() string { return e.message }

// Implement the ae.ErrorMessage interface
func (e CustomError) Message() string { return e.message }

// Implement the ae.ErrorCode interface
func (e CustomError) Code() string { return e.code }

// Use with ae functions
customErr := CustomError{"custom error", "CUSTOM_001"}
code := ae.Code(customErr)  // Returns "CUSTOM_001"
```
