package examples

import "github.com/aledantee/ae"

func Basic() {
	err := ae.New().
		Msg("something went wrong").
		Build()
	panic(err)
}
