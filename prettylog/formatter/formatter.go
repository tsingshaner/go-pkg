package formatter

import (
	"strings"
	"time"

	"github.com/tsingshaner/go-pkg/color"
	"github.com/tsingshaner/go-pkg/log/helper"
)

type Data = map[string]any

type Group struct {
	Key   string
	Value Data
}

type Log interface {
	Name() string
	Level() string
	Time() time.Time
	Msg() string
	Pid() int
	Src() string
	Err() string
	Stack() string
	Groups() []Group
}

var propertyPrefix = color.Bold(color.UnsafeDim("    » "))

func Formatter(log Log) string {
	switch {
	case strings.HasSuffix(log.Name(), helper.NameGinRouterLoggerSuffix):
		return FormatGinRouter(log)
	case strings.HasSuffix(log.Name(), helper.NameGORMLoggerSuffix):
		return FormatGorm(log)
	}

	sb := strings.Builder{}
	sb.WriteString(Level(log.Level()))

	if log.Pid() != 0 {
		sb.WriteByte(' ')
		sb.WriteString(Pid(log.Pid()))
	}

	sb.WriteByte(' ')
	sb.WriteString(Time(log.Time()))

	if log.Name() != "" {
		sb.WriteString(" #")
		sb.WriteString(color.UnsafeBold(color.UnsafeMagenta(log.Name())))
	}

	if log.Msg() != "" {
		sb.WriteString(color.UnsafeBold(color.UnsafeGreen(" ")))
		sb.WriteString(log.Msg())
		sb.WriteString(color.UnsafeGreen(""))
	}

	if log.Src() != "" {
		sb.WriteByte('\n')
		sb.WriteString(propertyPrefix)
		sb.WriteString(color.UnsafeBold(color.UnsafeBlue("src ")))
		sb.WriteString(color.UnsafeItalic(log.Src()))
	}

	if len(log.Groups()) > 0 {
		sb.WriteString(Groups(log.Groups(), propertyPrefix, "ctx").String())
	}

	sb.WriteByte('\n')

	if log.Stack() != "" {
		sb.WriteString(propertyPrefix)
		sb.WriteString(color.UnsafeBold(color.UnsafeBlue("stack ")))
		sb.WriteString(color.UnsafeDim(log.Stack()))
		sb.WriteByte('\n')
	}

	return sb.String()
}
