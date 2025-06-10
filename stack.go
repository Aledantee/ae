package ae

import "time"

// Stack represents a stack trace with associated metadata for a single goroutine.
type Stack struct {
	// ID is a unique identifier for this stack trace
	ID int `json:"id"`
	// State represents the current state of the stack (e.g. "running", "blocked")
	State string `json:"state"`
	// Wait indicates how long the stack has been waiting, if applicable
	Wait time.Duration `json:"wait"`
	// Locked indicates whether the stack is currently holding any locks
	Locked bool `json:"locked"`
	// Frames contains the stack frames in order from top to bottom
	Frames []*StackFrame `json:"frames"`
	// FramesElided indicates whether some frames were omitted from the trace
	FramesElided bool `json:"frames_elided"`
	// Parent points to the parent stack if this is a child stack
	Parent *Stack `json:"parent"`
	// Ancestor points to the root ancestor stack
	Ancestor *Stack `json:"ancestor"`
}

// StackFrame represents a single frame in a stack trace.
type StackFrame struct {
	// Func is the name of the function being called
	Func string `json:"func"`
	// File is the path to the source file
	File string `json:"file"`
	// Line is the line number in the source file
	Line int `json:"line"`
}

func newStack() []*Stack {
	return nil
}
