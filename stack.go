package ae

import (
	"bytes"
	"maps"
	"runtime/debug"
	"slices"
	"time"

	"github.com/DataDog/gostackparse"
)

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
	// CreatedBy points to the exact frame that created this stack.
	CreatedBy *StackFrame `json:"parent"`
	// Ancestor points to the root ancestor, which is the stack that crated this stack.
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

// newStack captures the current stack trace of all goroutines and returns them as a slice of Stack objects.
// It parses the debug stack information to extract goroutine details including their state, wait times,
// locked status, and stack frames. The function also establishes relationships between goroutines
// by linking them to their creating frames and ancestor stacks.
//
// Returns a slice of Stack objects representing all active goroutines.
func newStack() []*Stack {
	goRoutines, _ := gostackparse.Parse(bytes.NewReader(debug.Stack()))

	stacks := make(map[int]*Stack)
	ancestors := make(map[int]int)

	for _, g := range goRoutines {
		var frames []*StackFrame
		for _, frame := range g.Stack {
			frames = append(frames, &StackFrame{
				Func: frame.Func,
				File: frame.File,
				Line: frame.Line,
			})
		}

		stack := &Stack{
			ID:        g.ID,
			State:     g.State,
			Wait:      g.Wait,
			Locked:    g.LockedToThread,
			Frames:    frames,
			CreatedBy: nil,
			Ancestor:  nil,
		}

		if g.CreatedBy != nil {
			stack.CreatedBy = &StackFrame{
				Func: g.CreatedBy.Func,
				File: g.CreatedBy.File,
				Line: g.CreatedBy.Line,
			}
		}
		if g.Ancestor != nil {
			ancestors[g.ID] = g.Ancestor.ID
		}

		stacks[g.ID] = stack
	}

	for stackID, ancestorID := range ancestors {
		stack, ok := stacks[stackID]
		if !ok {
			// This is a bug, but let's not panic a prod system in the error path
			continue
		}

		ancestorStack, ok := stacks[ancestorID]
		if !ok {
			// This is a bug, but let's not panic a prod system in the error path
			continue
		}

		stack.Ancestor = ancestorStack
		stacks[stackID] = stack
	}

	return slices.Collect(maps.Values(stacks))
}
