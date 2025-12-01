package main

import "go.aledante.io/ae"

func main() {
	ae.Print(a())
}

func a() error {
	return func() error {
		return b()
	}()
}

func b() error {
	return func() error {
		return ae.New().
			Stack().
			Msg("error with stack")
	}()
}
