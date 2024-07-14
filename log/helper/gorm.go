package helper

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/tsingshaner/go-pkg/log"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
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
	opts.logger = l.Named(NameGORMLoggerSuffix)
	return &logger{opts}
}

// LogMode log mode
func (l *logger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	newLogger.LogLevel = slog.Level(level)
	return &newLogger
}

// Info print info
func (l *logger) Info(ctx context.Context, msg string, data ...any) {
	l.logger.Info(msg, slog.String("pos", utils.FileWithLineNum()), slog.Any("data", data))
}

// Warn print warn messages
func (l *logger) Warn(ctx context.Context, msg string, data ...any) {
	l.logger.Warn(msg, slog.String("pos", utils.FileWithLineNum()), slog.Any("data", data))
}

// Error print error messages
func (l *logger) Error(ctx context.Context, msg string, data ...any) {
	l.logger.Error(msg, slog.String("pos", utils.FileWithLineNum()), slog.Any("data", data))
}

func (l *logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if !l.logger.Enabled(log.SlogLevelTrace) {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.logger.Enabled(log.SlogLevelError) &&
		(!errors.Is(err, gormLogger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()

		l.logger.Error(err.Error(),
			slog.String("pos", utils.FileWithLineNum()),
			slog.String("sql", sql),
			slog.Float64("elapsed", float64(elapsed.Nanoseconds())/1e6),
			slog.Int64("rows", rows),
		)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.logger.Enabled(log.SlogLevelWarn):
		sql, rows := fc()

		l.logger.Warn(fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold),
			slog.String("pos", utils.FileWithLineNum()),
			slog.String("sql", sql),
			slog.Float64("elapsed", float64(elapsed.Nanoseconds())/1e6),
			slog.Int64("rows", rows),
		)
	case l.logger.Enabled(log.SlogLevelInfo):
		sql, rows := fc()

		l.logger.Info("",
			slog.String("pos", utils.FileWithLineNum()),
			slog.String("sql", sql),
			slog.Float64("elapsed", float64(elapsed.Nanoseconds())/1e6),
			slog.Int64("rows", rows),
		)
	}
}

// ParamsFilter filter params
func (l *logger) ParamsFilter(ctx context.Context, sql string, params ...any) (string, []any) {
	if l.GORMLoggerOptions.ParameterizedQueries {
		return sql, nil
	}
	return sql, params
}
