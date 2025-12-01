package main

import "go.aledante.io/ae"

func main() {
	err := ae.Wrap("test", ae.Msg("cause"))
	ae.Print(err)
}
