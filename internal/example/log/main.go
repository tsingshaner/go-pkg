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
	logger := log.NewSlog(w, &log.SlogHandlerOptions{
		Level:     log.LevelAll,
		AddSource: true,
	})

	child := logger.Child(
		slog.String("app", "custom_slog"),
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
}
