package ae

// ErrorExitCode defines an interface for errors that can provide an exit code.
type ErrorExitCode interface {
	// ErrorExitCode returns the exit code associated with the error.
	// If the error does not have an associated exit code, the highest exit code of all recursive causes is returned.
	ErrorExitCode() int
}

// ExitCode extracts the process exit code from an error.
// If the error implements ErrorExitCode, returns its ExitCode().
// Otherwise, recursively checks all causes and returns the highest exit code found.
// Returns 0 if err is nil, otherwise defaults to 1.
func ExitCode(err error) int {
	if err == nil {
		return 0
	}

	if ae, ok := err.(ErrorExitCode); ok && ae.ErrorExitCode() > 0 {
		return ae.ErrorExitCode()
	}

	exitCode := 1
	for _, cause := range Causes(err) {
		if ec := ExitCode(cause); ec > exitCode {
			exitCode = ec
		}
	}

	return exitCode
}
