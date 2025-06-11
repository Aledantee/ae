package main

import (
	"errors"
	"fmt"
	"github.com/aledantee/ae"
)

// This example demonstrates error chaining using ae.New() and its chaining methods.
func main() {
	err := ae.New().
		Cause(errors.New("invalid input")).
		Tag("validation").
		Msg("failed to process request")

	fmt.Println(err)
	// Output: failed to process request: invalid input {validation}
}
