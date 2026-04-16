// Example extract shows how to read metadata back out of an ae error using
// the package's accessor functions. Each accessor is interface-based, so it
// also works on any error that implements the corresponding ErrorXxx
// interface — not just values built with ae.New().
package main

import (
	"fmt"

	"go.aledante.io/ae"
)

func main() {
	err := ae.New().
		Code("E_IO").
		ExitCode(2).
		Tag("filesystem").
		Tag("transient").
		Attr("path", "/var/log/app.log").
		Attr("bytes", 4096).
		Hint("check disk space").
		TraceId("0123456789abcdef0123456789abcdef").
		SpanId("abcdef0123456789").
		Cause(ae.Msg("no space left on device")).
		Related(ae.Msg("metrics flush failed")).
		Now().
		UserMsg("write failed", "We couldn't save your changes. Please retry.")

	fmt.Println("Message      :", ae.Message(err))
	fmt.Println("UserMessage  :", ae.UserMessage(err))
	fmt.Println("Hint         :", ae.Hint(err))
	fmt.Println("Code         :", ae.Code(err))
	fmt.Println("ExitCode     :", ae.ExitCode(err))
	fmt.Println("TraceId      :", ae.TraceId(err))
	fmt.Println("SpanId       :", ae.SpanId(err))
	fmt.Println("Tags         :", ae.Tags(err))
	fmt.Println("Attributes   :", ae.Attributes(err))
	fmt.Println("Timestamp    :", ae.Timestamp(err).Format("2006-01-02T15:04:05Z07:00"))
	fmt.Println("Causes       :", ae.Causes(err))
	fmt.Println("Related      :", ae.Related(err))
	fmt.Println("Stacks       :", ae.Stacks(err))
}
