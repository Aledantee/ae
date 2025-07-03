package main

import "github.com/aledantee/ae"

func main() {
	err := ae.Wrap("test", ae.Msg("cause"))
	ae.Print(err)
}
