package formatter

import (
	"strings"

	"github.com/tsingshaner/go-pkg/color"
)

var (
	TagTrace = color.Bold(color.Cyan("trace"))
	TagDebug = color.Bold(color.Blue("debug"))
	TagInfo  = color.Bold(color.Green(" info"))
	TagWarn  = color.Bold(color.Yellow(" warn"))
	TagError = color.Bold(color.Magenta("error"))
	TagFatal = color.Bold(color.Red("fatal"))
	TagPanic = color.Bold(color.Red("panic"))
)

func Level(level string) string {
	l := strings.Trim(strings.ToLower(level), " ")

	switch l {
	case "trace":
		return TagTrace
	case "info":
		return TagInfo
	case "warn":
		return TagWarn
	case "error":
		return TagError
	case "debug":
		return TagDebug
	case "fatal":
		return TagFatal
	case "panic":
		return TagPanic
	}

	return color.Bold(color.Blue(l))
}
