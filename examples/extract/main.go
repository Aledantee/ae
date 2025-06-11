package main

import (
	"errors"
	"fmt"
	"github.com/aledantee/ae"
)

// This example demonstrates extracting an error message using ae.Message().
// It creates a wrapped error with fmt.Errorf and extracts its message.
func main() {
	err := fmt.Errorf("database connection failed: %w", errors.New("timeout"))
	message := ae.Message(err)
	fmt.Println(message)
	// Output: database connection failed: timeout
}
