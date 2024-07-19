package main

import (
	"os"
	"os/signal"

	"github.com/tsingshaner/go-pkg/log/console"
	"github.com/tsingshaner/go-pkg/prettylog"
	"github.com/tsingshaner/go-pkg/prettylog/adapter"
	"github.com/tsingshaner/go-pkg/prettylog/formatter"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		console.Info("waiting main program to exit or print ^C again to force exit")
	}()

	prettylog.JSONReader(func(ro *prettylog.ReaderOptions) {
		ro.Formatter = func(d formatter.Data, _ []byte) string {
			return formatter.Formatter(adapter.DefaultAdaptor(d, nil))
		}
	})
}
