package runtime

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStackTrace(t *testing.T) {
	caller, stack, err := GetStackTrace(true, true, 0)

	assert.Nil(t, err)

	assert.Equal(t, caller.Function, "github.com/tsingshaner/go-pkg/util/runtime.TestStackTrace")
	assert.Contains(t, caller.File, "util/runtime/stacktrace_test.go")
	assert.Greater(t, caller.Line, 5)

	assert.Contains(t, stack, "runtime.TestStackTrace")
	assert.Contains(t, stack, "testing.tRunner")
	assert.Equal(t, len(strings.Split(stack, "\n")), 4)
}

func TestStackTraceWithoutStack(t *testing.T) {
	caller, stack, err := GetStackTrace(true, false, 0)

	assert.Nil(t, err)

	assert.Equal(t, caller.Function, "github.com/tsingshaner/go-pkg/util/runtime.TestStackTraceWithoutStack")
	assert.Contains(t, caller.File, "util/runtime/stacktrace_test.go")
	assert.Greater(t, caller.Line, 5)

	assert.Equal(t, stack, "")
}

func TestStackTraceWithoutCaller(t *testing.T) {
	caller, stack, err := GetStackTrace(false, true, 0)

	assert.Nil(t, err)

	assert.Nil(t, caller)

	assert.Contains(t, stack, "runtime.TestStackTraceWithoutCaller")
	assert.Contains(t, stack, "testing.tRunner")
	assert.Equal(t, len(strings.Split(stack, "\n")), 4)
}
