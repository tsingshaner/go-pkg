package main

import (
	"io"
	"log/slog"
	"os"

	"github.com/tsingshaner/go-pkg/log"
)

func main() {
	w := createWriter()
	println("------custom slog------")
	exampleCustomSlog(w)
}

func createWriter() io.Writer {
	fw, err := log.NewFileWriter(func(config *log.FileConfig) {
		config.Filepath = "testdata/app.log"
		config.Compress = true
	})

	if err != nil {
		panic(err)
	}

	return log.NewWriter(fw, os.Stdout)
}

func exampleCustomSlog(w io.Writer) {
	logger, levelToggler := log.NewSlog(w, &log.SlogHandlerOptions{
		Level: slog.Level(log.LevelAll),
	}, func(o *log.Options) {
		o.AddSource = true
	})

	child := logger.Named("custom").Named("slog").Child(
		slog.String("version", "v1.0.0"),
		slog.Int("pid", os.Getpid()),
	)

	logger.Trace("custom slog trace")
	child.Trace("custom slog child trace")

	logger.Debug("custom slog debug")
	child.Debug("custom slog child debug")

	logger.Info("custom slog info")
	child.Info("custom slog child info")

	logger.Warn("custom slog warn")
	child.Warn("custom slog child warn")

	logger.Error("custom slog error")
	child.Error("custom slog child error")

	logger.Fatal("custom slog fatal")
	child.Fatal("custom slog child fatal")

	namedLogger := logger.Named("app").Named("user").Named("repo").Named("sql")
	namedLogger.Info("custom slog named logger")

	grouped := logger.WithGroup("group")
	grouped.Info("custom slog with group")
	grouped.Info("custom slog with group", slog.String("key", "value"))

	levelToggler(log.LevelError | log.LevelFatal)
	grouped.Debug("not print")
	grouped.Info("not print")
	grouped.Warn("not print")
	grouped.Error("error print")
	grouped.Fatal("fatal print", slog.String("key", "value"))

	childWithStack := grouped.WithOptions(&log.ChildLoggerOptions{
		AddSource:  true,
		StackTrace: log.LevelError | log.LevelFatal,
	})

	childWithStack.Debug("not print")
	childWithStack.Info("not print")
	childWithStack.Warn("not print")
	childWithStack.Error("error print")
	childWithStack.Fatal("fatal print", slog.String("key", "value"))
}
