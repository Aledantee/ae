// Example errors demonstrates the ae/errors sub-package: drop-in replacements
// for the stdlib errors package that produce ae.Ae values, plus the
// interop guarantees with errors.Is and errors.As.
package main

import (
	"fmt"

	"go.aledante.io/ae"
	aeerrors "go.aledante.io/ae/errors"
)

// sentinel is a sentinel error used to demonstrate errors.Is matching.
var sentinel = aeerrors.New("sentinel: not found")

func main() {
	// errors.New — creates an ae error, not a stdlib error.
	basic := aeerrors.New("plain error")
	fmt.Println("--- errors.New ---")
	ae.Print(basic)

	// errors.Join — collapses multiple errors into a single ae.Ae whose
	// causes are the inputs.
	joined := aeerrors.Join(
		aeerrors.New("first failure"),
		aeerrors.New("second failure"),
		nil, // nils are filtered out
		aeerrors.New("third failure"),
	)
	fmt.Println("\n--- errors.Join ---")
	ae.Print(joined)

	// errors.Is — walks the cause chain looking for a match. The sentinel
	// here lives two levels deep.
	wrapped := ae.Wrap("while loading config",
		ae.Wrap("while reading file", sentinel))
	fmt.Println("\n--- errors.Is ---")
	fmt.Println("matches sentinel?", aeerrors.Is(wrapped, sentinel))

	// errors.As — extracts a concrete type from anywhere in the chain.
	var target *ae.Ae
	if aeerrors.As(wrapped, &target) {
		fmt.Println("\n--- errors.As ---")
		fmt.Println("extracted *ae.Ae with message:", ae.Message(target))
	}
}
