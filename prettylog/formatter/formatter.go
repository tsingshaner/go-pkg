package formatter

import (
	"strings"
	"time"

	"github.com/tsingshaner/go-pkg/color"
)

type Data map[string]any

type Log interface {
	Level() string
	Timestamp() time.Time
	Msg() string
	Pid() int
	Src() string
	Err() string
	Data() Data
}

func Formatter(log Log) string {
	sb := strings.Builder{}
	sb.WriteString(Level(log.Level()))

	if log.Pid() != 0 {
		sb.WriteString(" ")
		sb.WriteString(Pid(log.Pid()))
	}

	sb.WriteString(" ")
	sb.WriteString(Time(log.Timestamp()))

	if log.Msg() != "" {
		sb.WriteString(color.UnsafeBold(color.UnsafeGreen(" ")))
		sb.WriteString(log.Msg())
		sb.WriteString(color.UnsafeGreen(""))
	}

	if log.Src() != "" {
		sb.WriteString(color.UnsafeBold(color.UnsafeGreen("\n # ")))
		sb.WriteString(color.UnsafeItalic(log.Src()))
		sb.WriteString(color.UnsafeItalic(color.UnsafeCyan("()")))
	}

	if len(log.Data()) > 0 {
		sb.WriteString("\n")
		sb.WriteString(Map(log.Data(), color.Bold(color.UnsafeDim(" » ")), "ctx", 2).String())
	}
	sb.WriteString("\n")

	return sb.String()
}
