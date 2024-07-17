package helper

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/tsingshaner/go-pkg/log"
	gormLogger "gorm.io/gorm/logger"
)

const NameGORMLoggerSuffix = "gorm.__gorm"

// Config logger config
type GORMLoggerOptions struct {
	SlowThreshold             time.Duration
	UseConsole                bool
	IgnoreRecordNotFoundError bool
	ParameterizedQueries      bool
	LogLevel                  slog.Level
	logger                    log.Slog
}

type logger struct {
	GORMLoggerOptions
}

func NewGormLogger(l log.Slog, opts GORMLoggerOptions) gormLogger.Interface {
	opts.logger = l.Named(NameGORMLoggerSuffix).WithOptions(&log.ChildLoggerOptions{
		AddSource:  true,
		SkipCaller: 2,
	})
	return &logger{opts}
}

// LogMode log mode
func (l *logger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	newLogger.LogLevel = slog.Level(level)
	return &newLogger
}

// Info print info
func (l *logger) Info(_ context.Context, msg string, data ...any) {
	l.logger.Info(fmt.Sprintf(msg, data...))
}

// Warn print warn messages
func (l *logger) Warn(_ context.Context, msg string, data ...any) {
	l.logger.Warn(fmt.Sprintf(msg, data...))
}

// Error print error messages
func (l *logger) Error(_ context.Context, msg string, data ...any) {
	l.logger.Error(fmt.Sprintf(msg, data...))
}

type traceRecord struct {
	elapsed       time.Duration
	getSqlAndRows func() (string, int64)
}

func (tr *traceRecord) Attrs() []slog.Attr {
	sql, rows := tr.getSqlAndRows()

	return []slog.Attr{
		slog.String("sql", sql),
		slog.Float64("elapsed", float64(tr.elapsed.Nanoseconds())/1e6),
		slog.Int64("rows", rows),
	}
}

func (tr *traceRecord) IsSlow(slowThreshold time.Duration) bool {
	return tr.elapsed > slowThreshold && slowThreshold != 0
}

func (l *logger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if !l.logger.Enabled(log.SlogLevelTrace) {
		return
	}

	tr := &traceRecord{time.Since(begin), fc}
	switch {
	case l.shouldLogAsError(err):
		l.logger.Error(err.Error(), tr.Attrs()...)
	case l.logger.Enabled(log.SlogLevelWarn) && tr.IsSlow(l.SlowThreshold):
		l.logger.Warn(fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold), tr.Attrs()...)
	case l.logger.Enabled(log.SlogLevelInfo):
		l.logger.Info("", tr.Attrs()...)
	}
}

func (l *logger) shouldLogAsError(e error) bool {
	if e != nil && l.logger.Enabled(log.SlogLevelError) {
		return !l.IgnoreRecordNotFoundError || !errors.Is(e, gormLogger.ErrRecordNotFound)
	}

	return false
}

// ParamsFilter filter params
func (l *logger) ParamsFilter(ctx context.Context, sql string, params ...any) (string, []any) {
	if l.GORMLoggerOptions.ParameterizedQueries {
		return sql, nil
	}
	return sql, params
}
