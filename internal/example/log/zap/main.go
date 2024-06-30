package main

import (
	"io"
	"os"

	"github.com/tsingshaner/go-pkg/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	w := createWriter()
	println("------custom zap------")
	exampleCustomZap(w)
}

func createWriter() io.Writer {
	fw, err := log.NewFileWriter(func(config *log.FileConfig) {
		config.Filepath = "testdata/log_zap.log"
		config.Compress = true
	})

	if err != nil {
		panic(err)
	}

	return log.NewWriter(fw, os.Stdout)
}

func exampleCustomZap(w io.Writer) {
	logger := log.NewZapLogger(log.NewZapCore(
		log.NewZapJSONEncoder(),
		zapcore.AddSync(os.Stdout),
		log.LevelAll,
	), zap.AddCaller(), zap.AddCallerSkip(2))

	child := logger.Named("custom").Named("zap").Child(
		zap.String("version", "v1.0.0"),
		zap.Int("pid", os.Getpid()),
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

	namedLogger := logger.Named("app").Named("user").Named("repo").Named("sql")
	namedLogger.Info("custom zap named logger")

	grouped := logger.WithGroup("group")
	grouped.Info("custom zap with group")
	grouped.Info("custom zap with group", zap.String("key", "value"))

	logger.Fatal("custom zap fatal")
	child.Fatal("custom zap child fatal")
}
