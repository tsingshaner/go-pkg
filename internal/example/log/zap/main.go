package main

import (
	"log/slog"
	"os"

	"github.com/tsingshaner/go-pkg/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	println("------custom zap------")
	exampleCustomZap()
}

func exampleCustomZap() {
	core, levelToggler := log.NewZapCore(
		log.NewZapJSONEncoder(),
		zapcore.AddSync(os.Stdout),
		log.LevelAll,
	)

	stackLevelFunc, stackLevelToggler := log.NewZapLevelFilter(log.LevelError | log.LevelFatal)

	logger := log.NewZapLog(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(2),
		zap.AddStacktrace(stackLevelFunc),
	)

	child := logger.Named("custom").Named("zap").Child(
		slog.String("version", "v1.0.0"),
		slog.Int("pid", os.Getpid()),
	)

	logger.Trace("custom zap trace")
	child.Trace("custom zap child trace")

	logger.Debug("custom zap debug")
	child.Debug("custom zap child debug")

	logger.Info("custom zap info")
	child.Info("custom zap child info")

	logger.Warn("custom zap warn")
	child.Warn("custom zap child warn")

	logger.Error("custom zap error")
	child.Error("custom zap child error")

	logger.Fatal("custom zap fatal")
	child.Fatal("custom zap child fatal")

	namedLogger := logger.Named("app").Named("user").Named("repo").Named("sql")
	namedLogger.Info("custom zap named logger")

	grouped := logger.WithGroup("group")
	grouped.Info("custom zap with group")
	grouped.Info("custom zap with group", slog.String("key", "value"))

	levelToggler(log.LevelError | log.LevelFatal)
	namedLogger.Trace("not print")
	namedLogger.Debug("not print")
	namedLogger.Info("not print")
	namedLogger.Warn("not print")
	namedLogger.Error("error print")
	namedLogger.Fatal("fatal print", slog.String("err", "fatal error"))

	stackLevelToggler(log.LevelFatal)
	namedLogger.Error("error print without stack")
	namedLogger.Fatal("fatal print with stack")
}
