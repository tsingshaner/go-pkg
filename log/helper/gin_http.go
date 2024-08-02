package helper

import (
	"bytes"
	"log/slog"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tsingshaner/go-pkg/log"
)

const NameGinHttpLoggerSuffix = "gin.__http_logger"

type Options struct {
	TraceIDExtractor func(*gin.Context) string
	Logger           log.Slog
}

var MaskFieldRegs = []*regexp.Regexp{
	regexp.MustCompile(`("pwd"|"password"|"token"|"accessToken"|"refreshToken")\s*:\s*"[^"]*"`),
}

func maskSecrets(data string) string {
	maskedData := data
	for _, re := range MaskFieldRegs {
		maskedData = re.ReplaceAllString(maskedData, `$1:"***"$2`)
	}

	return maskedData
}

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func New(opts *Options) gin.HandlerFunc {
	logger := opts.Logger.WithOptions(&log.ChildLoggerOptions{
		AddSource:  false,
		StackTrace: log.LevelSilent,
	}).Named(NameGinHttpLoggerSuffix)

	return func(c *gin.Context) {
		begin := time.Now()
		body := &bytes.Buffer{}
		c.Writer = &CustomResponseWriter{c.Writer, body}

		c.Next()

		if !logger.Enabled(slog.Level(log.LevelTrace)) {
			return
		}

		req := make([]any, 0, 11)
		req = append(req,
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("query", c.Request.URL.RawQuery),
			slog.String("clientIP", c.ClientIP()),
			slog.String("userAgent", c.Request.UserAgent()),
			slog.String("referer", c.Request.Referer()),
			slog.String("proto", c.Request.Proto),
			slog.String("host", c.Request.Host),
			slog.String("contentType", c.Request.Header.Get("Content-Type")),
			slog.String("remoteAddr", c.Request.RemoteAddr),
		)

		if body, err := c.GetRawData(); err == nil {
			req = append(req, slog.String("body", maskSecrets(string(body))))
		}

		reqGroup := slog.Group("req", req...)

		logger.Trace(http.StatusText(c.Writer.Status()),
			slog.String("traceID", opts.TraceIDExtractor(c)),
			slog.String("latency", time.Since(begin).String()),
			reqGroup,

			slog.Group("resp",
				slog.String("body", maskSecrets(body.String())),
				slog.Int("status", c.Writer.Status()),
				slog.Int("size", c.Writer.Size()),
				slog.String("contentType", c.Writer.Header().Get("Content-Type")),
			))
	}
}
