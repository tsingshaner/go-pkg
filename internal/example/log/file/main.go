package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/tsingshaner/go-pkg/log"
	"github.com/tsingshaner/go-pkg/log/console"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	writer, err := log.NewFileWriter(func(c *log.FileConfig) {
		c.Filepath = "testdata/rotate.log"
		c.BackupTime = 24 * 7
		// c.TimeClock
		c.Compress = true
		c.BackupNum = 3
		c.MaxSize = 1024
	})

	if err != nil {
		console.Fatal("create file writer failed: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log.NewFilesClear(ctx, func(c *log.FileClearConfig) {
		c.BackupTime = 3
		c.TimeUnit = time.Second
		c.BackupNum = 3
		c.Compress = true
		c.CheckInterval = 5 * time.Second
	})

	core, _ := log.NewZapCore(log.NewZapJSONEncoder(), zapcore.AddSync(writer), log.LevelAll)
	logger := log.NewZapLog(core, zap.AddCaller(), zap.AddCallerSkip(2))

	for i := 0; i < 100; i++ {
		logger.Trace("hello", slog.String("name", "jack"))
		logger.Debug("hello", slog.String("name", "jack"))
		logger.Info("hello", slog.String("name", "jack"))
		logger.Warn("hello", slog.String("name", "jack"))
		logger.Error("hello", slog.String("name", "jack"))
		logger.Fatal("hello", slog.String("name", "jack"))

		time.Sleep(time.Second)
	}
}
