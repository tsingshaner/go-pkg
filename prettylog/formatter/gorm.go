package formatter

import (
	"strconv"
	"strings"
	"time"

	"github.com/tsingshaner/go-pkg/color"
)

func FormatGorm(log Log) string {
	gormLog, _ := parseGormLog(getLastGroup(log).Value)

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

	if gormLog.rows != 0 {
		sb.WriteString(" ")
		sb.WriteString(color.UnsafeGray("rows="))
		sb.WriteString(color.UnsafeCyan(strconv.Itoa(gormLog.rows)))
	}

	if gormLog.duration != 0 {
		sb.WriteString(" ")
		sb.WriteString(color.UnsafeItalic(color.UnsafeYellow(gormLog.duration.String())))
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

	if log.Stack() != "" {
		sb.WriteByte('\n')
		sb.WriteString(propertyPrefix)
		sb.WriteString(color.UnsafeBold(color.UnsafeBlue("stack ")))
		sb.WriteString(color.UnsafeDim(log.Stack()))
	}

	if gormLog.sql != "" {
		sb.WriteByte('\n')
		sb.WriteString(propertyPrefix)
		sb.WriteString(color.UnsafeBold(color.UnsafeBlue("sql ")))
		sb.WriteString(color.UnsafeGray(gormLog.sql))
	}

	sb.WriteByte('\n')
	return sb.String()
}

type gormLog struct {
	rows     int
	sql      string
	duration time.Duration
}

func parseGormLog(data Data) (*gormLog, Data) {
	log := &gormLog{}
	for k, v := range data {
		switch k {
		case "rows":
			if rows, ok := v.(float64); ok {
				log.rows = int(rows)
				delete(data, k)
			}
		case "sql":
			if sql, ok := v.(string); ok {
				log.sql = sql
				delete(data, k)
			}
		case "elapsed":
			if duration, ok := v.(float64); ok {
				log.duration = time.Duration(duration) * time.Millisecond
				delete(data, k)
			}
		}
	}
	return log, data
}
