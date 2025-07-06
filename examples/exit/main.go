package main

import (
	"github.com/aledantee/ae"
	"os"
)

func main() {
	err := ae.New().
		Cause(ae.New().ExitCode(200).Msg("cause")).
		Msg("should exit with 200")

	os.Exit(ae.ExitCode(err))
}
