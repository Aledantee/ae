package ae

// Message extracts the internal error message from an error.
// If the error implements ErrorMessage, returns its Message().
// Otherwise, returns the error's Error() string.
// Returns an empty string if err is nil.
func Message(err error) string {
	if err == nil {
		return ""
	}

	if ae, ok := err.(ErrorMessage); ok {
		return ae.Message()
	}

	return err.Error()
}

// UserMessage extracts the user-friendly error message from an error.
// If the error implements ErrorUserMessage, returns its UserMessage().
// Returns an empty string if err is nil or if the error does not implement ErrorUserMessage.
func UserMessage(err error) string {
	if err == nil {
		return ""
	}

	if ae, ok := err.(ErrorUserMessage); ok {
		return ae.UserMessage()
	}

	return ""
}

// Hint extracts the resolution hint from an error.
// If the error implements ErrorHint, returns its Hint().
// Returns an empty string if err is nil or if the error does not implement ErrorHint.
func Hint(err error) string {
	if err == nil {
		return ""
	}

	if ae, ok := err.(ErrorHint); ok {
		return ae.Hint()
	}

	return ""
}

// Code extracts the error code from an error.
// If the error implements ErrorCode, returns its Code().
// Returns an empty string if err is nil or if the error does not implement ErrorCode.
func Code(err error) string {
	if err == nil {
		return ""
	}

	if ae, ok := err.(ErrorCode); ok {
		return ae.Code()
	}

	return ""
}

// ExitCode extracts the process exit code from an error.
// If the error implements ErrorExitCode, returns its ExitCode().
// Otherwise, recursively checks all causes and returns the highest exit code found.
// Returns 1 if err is nil or if no exit code is found in the error or its causes.
func ExitCode(err error) int {
	if err == nil {
		return 1
	}

	if ae, ok := err.(ErrorExitCode); ok {
		return ae.ExitCode()
	}

	exitCode := 1
	for _, cause := range Causes(err) {
		if ec := ExitCode(cause); ec > exitCode {
			exitCode = ec
		}
	}

	return exitCode
}

// TraceId extracts the distributed tracing ID from an error.
// If the error implements ErrorTraceId, returns its TraceId().
// Returns an empty string if err is nil or if the error does not implement ErrorTraceId.
func TraceId(err error) string {
	if err == nil {
		return ""
	}

	if ae, ok := err.(ErrorTraceId); ok {
		return ae.TraceId()
	}

	return ""
}

// SpanId extracts the operation span ID from an error.
// If the error implements ErrorSpanId, returns its SpanId().
// Returns an empty string if err is nil or if the error does not implement ErrorSpanId.
func SpanId(err error) string {
	if err == nil {
		return ""
	}

	if ae, ok := err.(ErrorSpanId); ok {
		return ae.SpanId()
	}

	return ""
}

// Tags extracts the list of tags from an error.
// If the error implements ErrorTags, returns its Tags().
// Returns nil if err is nil or if the error does not implement ErrorTags.
func Tags(err error) []string {
	if err == nil {
		return nil
	}

	if ae, ok := err.(ErrorTags); ok {
		return ae.Tags()
	}

	return nil
}

// Attributes extracts the map of attributes from an error.
// If the error implements ErrorAttributes, returns its Attributes().
// Returns an empty map if err is nil or if the error does not implement ErrorAttributes.
func Attributes(err error) map[string]any {
	if err == nil {
		return make(map[string]any)
	}

	if ae, ok := err.(ErrorAttributes); ok {
		attrs := ae.Attributes()
		if attrs != nil {
			return attrs
		}
	}

	return make(map[string]any)
}

// Causes extracts the list of underlying causes from an error.
// If the error implements ErrorCauses, returns its Causes().
// If the error implements Unwrap() []error, returns its Unwrap().
// If the error implements Unwrap() error, returns a single-element slice containing its Unwrap().
// Returns nil if err is nil or if the error does not implement any of these interfaces.
func Causes(err error) []error {
	if err == nil {
		return nil
	}

	if ae, ok := err.(ErrorCauses); ok {
		return ae.Causes()
	}

	if ae, ok := err.(interface{ Unwrap() []error }); ok {
		return ae.Unwrap()
	}

	if ae, ok := err.(interface{ Unwrap() error }); ok {
		return []error{ae.Unwrap()}
	}

	return nil
}

// Related extracts the list of related errors from an error.
// If the error implements ErrorRelated, returns its Related().
// Returns nil if err is nil or if the error does not implement ErrorRelated.
func Related(err error) []error {
	if err == nil {
		return nil
	}

	if ae, ok := err.(ErrorRelated); ok {
		return ae.Related()
	}

	return nil
}

// Stacks extracts the list of stack traces from an error.
// If the error implements ErrorStacks, returns its Stacks().
// Returns nil if err is nil or if the error does not implement ErrorStacks.
func Stacks(err error) []*Stack {
	if err == nil {
		return nil
	}

	if ae, ok := err.(ErrorStacks); ok {
		return ae.Stacks()
	}

	return nil
}
