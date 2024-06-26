package main

import (
	"github.com/tsingshaner/go-pkg/prettylog"
	"github.com/tsingshaner/go-pkg/prettylog/adapter"
	"github.com/tsingshaner/go-pkg/prettylog/formatter"
)

func main() {
	prettylog.JSONReader(func(ro *prettylog.ReaderOptions) {
		ro.Formatter = func(d formatter.Data, b []byte) string {
			return formatter.Formatter(adapter.SlogAdaptor(d, b))
		}
	})
}
