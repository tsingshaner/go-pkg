package adapter

import (
	"time"

	"github.com/tsingshaner/go-pkg/prettylog/formatter"
)

type defaultLog struct {
	name   string
	level  string
	time   time.Time
	msg    string
	pid    int
	src    string
	err    string
	stack  string
	groups []formatter.Group
	data   *formatter.Node
}

func DefaultAdaptor(data formatter.Data, _ []byte) formatter.Log {
	l := &defaultLog{
		data:   parseNode("", data),
		groups: parseGroups("", data),
	}

	for _, group := range l.groups {
		if level, ok := group.Value["level"].(string); ok {
			l.level = level
			delete(group.Value, "level")
		}

		if name, ok := group.Value["name"].(string); ok {
			l.name = name
			delete(group.Value, "name")
		}

		if ts, ok := group.Value["time"].(string); ok {
			if date, err := time.Parse(time.RFC3339Nano, ts); err == nil {
				l.time = date
				delete(group.Value, "time")
			}
		}

		if msg, ok := group.Value["msg"].(string); ok {
			l.msg = msg
			delete(group.Value, "msg")
		}

		if pid, ok := group.Value["pid"].(float64); ok {
			l.pid = int(pid)
			delete(group.Value, "pid")
		}

		if src, ok := group.Value["src"].(string); ok {
			l.src = src

			delete(group.Value, "src")
		}

		if err, ok := group.Value["err"].(string); ok {
			l.err = err
			delete(group.Value, "err")
		}

		if stack, ok := group.Value["stack"].(string); ok {
			l.stack = stack
			delete(group.Value, "stack")
		}
	}

	return l
}

func (dl *defaultLog) Tree() *formatter.Node {
	return dl.data
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

func (dl *defaultLog) Groups() []formatter.Group {
	return dl.groups
}

func parseGroups(name string, data formatter.Data) []formatter.Group {
	groups := []formatter.Group{
		{Key: name, Value: data},
	}

	for k, v := range data {
		if sub, ok := v.(formatter.Data); ok {
			delete(data, k)
			groups = append(groups, parseGroups(k, sub)...)
			break
		}
	}

	return groups
}

func parseNode(name string, data formatter.Data) *formatter.Node {
	node := &formatter.Node{
		Key:  name,
		Data: data,
	}

	for k, v := range data {
		if sub, ok := v.(formatter.Data); ok {
			delete(data, k)
			node.Children = append(node.Children, parseNode(k, sub))
		}
	}

	return node
}
