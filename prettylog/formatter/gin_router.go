package formatter

import (
	"fmt"
	"strings"

	"github.com/tsingshaner/go-pkg/color"
)

func FormatGinRouter(log Log) string {
	sb := &strings.Builder{}

	if method, ok := log.Data()["method"].(string); ok {
		sb.WriteString(formatMethod(method))
		sb.WriteByte(' ')
	}

	if path, ok := log.Data()["path"].(string); ok {
		sb.WriteString(formatPath(path))
		sb.WriteByte(' ')
	}

	if layers, ok := log.Data()["layers"].(float64); ok {
		sb.WriteString(color.UnsafeGray("--{"))
		sb.WriteString(color.UnsafeCyan(fmt.Sprintf("%d", int(layers))))
		sb.WriteString(color.UnsafeGray("}-> "))
	}

	sb.WriteString(formatHandlerName(log.Msg()))

	return sb.String()
}

var (
	methodGet     = color.UnsafeGreen("    GET")
	methodPost    = color.UnsafeYellow("   POST")
	methodPut     = color.UnsafeBlue("    PUT")
	methodPatch   = color.UnsafeMagenta("  PATCH")
	methodOptions = color.UnsafeCyan("OPTIONS")
	methodDelete  = color.UnsafeRed(" DELETE")
)

func formatMethod(method string) string {
	switch strings.ToUpper(method) {
	case "GET":
		return methodGet
	case "POST":
		return methodPost
	case "PUT":
		return methodPut
	case "PATCH":
		return methodPatch
	case "OPTIONS":
		return methodOptions
	case "DELETE":
		return methodDelete
	default:
		return color.Bold(color.UnsafeGreen(method))
	}
}

func formatPath(path string) string {
	return color.Gray(strings.Join(
		strings.Split(fmt.Sprintf("%-40s", path), "/"),
		color.Black("/"),
	))
}

func formatHandlerName(handlerName string) string {
	handlerNameTokens := strings.Split(handlerName, "/")
	return color.UnsafeGray(strings.Join(handlerNameTokens[len(handlerNameTokens)-2:], "/"))
}
