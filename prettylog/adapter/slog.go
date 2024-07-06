package adapter

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/tsingshaner/go-pkg/color"
	"github.com/tsingshaner/go-pkg/prettylog/formatter"
)

type slogLog struct {
	name  string
	level string
	time  time.Time
	msg   string
	pid   int
	src   *slog.Source
	err   string
	stack string
	data  formatter.Data
}

func SlogAdaptor(data formatter.Data, _ []byte) formatter.Log {
	l := &slogLog{
		data: data,
	}

	if name, ok := data["name"].(string); ok {
		l.name = name
		delete(l.data, "name")
	}

	if level, ok := data[slog.LevelKey].(string); ok {
		l.level = level
		delete(l.data, slog.LevelKey)
	}

	if ts, ok := data[slog.TimeKey].(string); ok {
		if date, err := time.Parse(time.RFC3339Nano, ts); err == nil {
			l.time = date
			delete(l.data, slog.TimeKey)
		}
	}

	if msg, ok := data[slog.MessageKey].(string); ok {
		l.msg = msg
		delete(l.data, slog.MessageKey)
	}

	if pid, ok := data["pid"].(float64); ok {
		l.pid = int(pid)
		delete(l.data, "pid")
	}

	if src, ok := data[slog.SourceKey].(map[string]any); ok {
		l.src = &slog.Source{
			Function: src["function"].(string),
			File:     src["file"].(string),
			Line:     int(src["line"].(float64)),
		}

		delete(l.data, slog.SourceKey)
	}

	if err, ok := data["err"].(string); ok {
		l.err = err
		delete(l.data, "err")
	}

	if stack, ok := data["stack"].(string); ok {
		l.stack = stack
		delete(l.data, "stack")
	}

	return l
}

func (sl *slogLog) Stack() string {
	return sl.stack
}

func (sl *slogLog) Name() string {
	return sl.name
}

func (sl *slogLog) Level() string {
	return sl.level
}

func (sl *slogLog) Time() time.Time {
	return sl.time
}

func (sl *slogLog) Msg() string {
	return sl.msg
}

func (sl *slogLog) Pid() int {
	return sl.pid
}

func (sl *slogLog) Src() string {
	sb := strings.Builder{}
	sb.WriteString(sl.src.File)
	sb.WriteString(color.UnsafeDim(":"))
	sb.WriteString(fmt.Sprintf("%d ", sl.src.Line))
	sb.WriteString(color.UnsafeItalic(sl.src.Function))

	return sb.String()
}

func (sl *slogLog) Err() string {
	return sl.err
}

func (sl *slogLog) Data() formatter.Data {
	return sl.data
}
