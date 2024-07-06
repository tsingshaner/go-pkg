package adapter

import (
	"time"

	"github.com/tsingshaner/go-pkg/prettylog/formatter"
)

type defaultLog struct {
	name  string
	level string
	time  time.Time
	msg   string
	pid   int
	src   string
	err   string
	stack string
	data  formatter.Data
}

func DefaultAdaptor(data formatter.Data, _ []byte) formatter.Log {
	l := &defaultLog{data: data}

	if level, ok := data["level"].(string); ok {
		l.level = level
		delete(l.data, "level")
	}

	if name, ok := data["name"].(string); ok {
		l.name = name
		delete(l.data, "name")
	}

	if ts, ok := data["time"].(string); ok {
		if date, err := time.Parse(time.RFC3339Nano, ts); err == nil {
			l.time = date
			delete(l.data, "time")
		}
	}

	if msg, ok := data["msg"].(string); ok {
		l.msg = msg
		delete(l.data, "msg")
	}

	if pid, ok := data["pid"].(float64); ok {
		l.pid = int(pid)
		delete(l.data, "pid")
	}

	if src, ok := data["src"].(string); ok {
		l.src = src

		delete(l.data, "src")
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

func (dl *defaultLog) Name() string {
	return dl.name
}

func (dl *defaultLog) Level() string {
	return dl.level
}

func (dl *defaultLog) Time() time.Time {
	return dl.time
}

func (dl *defaultLog) Msg() string {
	return dl.msg
}

func (dl *defaultLog) Pid() int {
	return dl.pid
}

func (dl *defaultLog) Src() string {
	return dl.src
}

func (dl *defaultLog) Err() string {
	return dl.err
}

func (dl *defaultLog) Stack() string {
	return dl.stack
}

func (dl *defaultLog) Data() formatter.Data {
	return dl.data
}
