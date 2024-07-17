package console

import (
	"fmt"
	"os"
	"strings"

	"github.com/tsingshaner/go-pkg/log"
	"github.com/tsingshaner/go-pkg/prettylog/formatter"
	"github.com/tsingshaner/go-pkg/util/bitmask"
)

func init() {
	SetLevel(log.LevelAll)
}

var (
	level = log.LevelAll
	Trace = noop
	Debug = noop
	Info  = noop
	Warn  = noop
	Error = noop
	Fatal = noop
)

func noop(_ string, _ ...any) {}

func trace(format string, s ...any) {
	fmt.Printf(buildFormat(formatter.TagTrace, format).String(), s...)
}

func debug(format string, s ...any) {
	fmt.Printf(buildFormat(formatter.TagDebug, format).String(), s...)
}

func info(format string, s ...any) {
	fmt.Printf(buildFormat(formatter.TagInfo, format).String(), s...)
}

func warn(format string, s ...any) {
	fmt.Printf(buildFormat(formatter.TagWarn, format).String(), s...)
}

func error(format string, s ...any) {
	fmt.Printf(buildFormat(formatter.TagError, format).String(), s...)
}

func fatal(format string, s ...any) {
	fmt.Printf(buildFormat(formatter.TagFatal, format).String(), s...)
	os.Exit(1)
}

func SetLevel(l log.Level) {
	level = l

	if bitmask.Has(level, log.LevelTrace) {
		Trace = trace
	} else {
		Trace = noop
	}

	if bitmask.Has(level, log.LevelDebug) {
		Debug = debug
	} else {
		Debug = noop
	}

	if bitmask.Has(level, log.LevelInfo) {
		Info = info
	} else {
		Info = noop
	}

	if bitmask.Has(level, log.LevelWarn) {
		Warn = warn
	} else {
		Warn = noop
	}

	if bitmask.Has(level, log.LevelError) {
		Error = error
	} else {
		Error = noop
	}

	if bitmask.Has(level, log.LevelFatal) {
		Fatal = fatal
	} else {
		Fatal = noop
	}
}

func buildFormat(prefix, format string) *strings.Builder {
	builder := &strings.Builder{}
	builder.WriteString(prefix)
	builder.WriteByte(' ')
	builder.WriteString(format)
	builder.WriteByte('\n')
	return builder
}
