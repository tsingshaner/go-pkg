package helper

import (
	"log/slog"

	"github.com/tsingshaner/go-pkg/log"
)

const NameGinRouterLoggerSuffix = "gin.__router"

type GinRouterLogger func(httpMethod, absolutePath, handlerName string, nuHandlers int)

func NewGinRouterLogger(logger log.Slog) GinRouterLogger {
	routerLog := logger.Named(NameGinRouterLoggerSuffix)

	return func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		routerLog.Debug(
			handlerName,
			slog.String("method", httpMethod),
			slog.String("path", absolutePath),
			slog.Int("layers", nuHandlers),
		)
	}
}
