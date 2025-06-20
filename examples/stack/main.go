package main

import "github.com/aledantee/ae"

func main() {
	ae.Print(a(), ae.WithJSON())
}

func a() error {
	return b()
}

func b() error {
	return ae.New().
		Stack().
		Msg("error with stack")
}
