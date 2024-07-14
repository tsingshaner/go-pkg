package runtime

import (
	"fmt"
)

type Caller struct {
	Function string
	File     string
	Line     int
}

// GetStackTrace return Caller and stack string, you can only use stack or caller
func GetStackTrace(addCaller, addStack bool, skip int) (*Caller, string, error) {
	// Adding the caller or stack trace requires capturing the callers of
	// this function. We'll share information between these two.
	stackDepth := First
	if addStack {
		stackDepth = Full
	}
	stack := Capture(skip+1, stackDepth)
	defer stack.Free()

	if stack.Count() == 0 && addCaller {
		if addCaller {
			return nil, "", fmt.Errorf("failed to get caller")
		}
		return nil, "", nil
	}

	frame, more := stack.Next()

	caller := (*Caller)(nil)
	if addCaller {
		caller = &Caller{frame.Function, frame.File, frame.Line}
	}

	if addStack {
		buffer := GetBufferPool()
		defer buffer.Free()

		formatter := NewFormatter(buffer)

		// We've already extracted the first frame, so format that
		// separately and defer to stackfmt for the rest.
		formatter.FormatFrame(frame)
		if more {
			formatter.FormatStack(stack)
		}
		return caller, buffer.String(), nil
	}

	return caller, "", nil
}
