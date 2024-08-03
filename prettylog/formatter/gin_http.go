package formatter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tsingshaner/go-pkg/color"
)

var (
	reqPrefix      = color.UnsafeGray("  «")
	respPrefix     = color.UnsafeBold(color.UnsafeGreen(" »»"))
	bodyField      = color.UnsafeBlue("Body")
	userAgentField = color.UnsafeMagenta("UA")

	leftSimpleQuote  = color.UnsafeGreen("[")
	rightSimpleQuote = color.UnsafeGreen("]")

	bodyLeftSimpleQuote  = color.UnsafeYellow("\"")
	bodyRightSimpleQuote = color.UnsafeYellow("\"")

	referrerField = color.UnsafeMagenta("Referrer")
)

func FormatGinHttp(log Log) string {
	sb := &strings.Builder{}

	var req, resp Data
	if log.Tree().Children[0].Key == "req" {
		req = log.Tree().Children[0].Data
		resp = log.Tree().Children[1].Data
	} else {
		req = log.Tree().Children[1].Data
		resp = log.Tree().Children[0].Data
	}

	if status, ok := resp["status"].(float64); ok {
		fmt.Fprintf(sb, "  %s", formatStatus(status))
	}

	if log.Pid() != 0 {
		fmt.Fprintf(sb, " %s", Pid(log.Pid()))
	}

	fmt.Fprintf(sb, " %s", Time(log.Time()))

	if traceID, ok := log.Tree().Data["traceID"].(string); ok && traceID != "" {
		fmt.Fprintf(sb, " %s", color.UnsafeCyan(traceID))
	}

	if uri := formatRequestUri(req); uri != "" {
		fmt.Fprintf(sb, " %s", color.UnsafeItalic(uri))
	}

	if latency, ok := log.Tree().Data["latency"].(string); ok {
		fmt.Fprintf(sb, " +%s", color.UnsafeYellow(latency))
	}

	sb.WriteByte('\n')

	if method, ok := req["method"].(string); ok {
		fmt.Fprintf(sb, "  %s %s", reqPrefix, strings.ReplaceAll(formatMethod(method), " ", ""))
	}

	if query, ok := req["query"].(string); ok && query != "" {
		fmt.Fprintf(sb, " %s", color.UnsafeBlue(query))
	}

	if contentType, ok := req["contentType"].(string); ok && contentType != "" {
		fmt.Fprintf(sb, " %s", formatContentType(contentType))
	}

	if proto, ok := req["proto"].(string); ok && proto != "" {
		fmt.Fprintf(sb, " %s", proto)
	}

	if remoteAddr, ok := req["remoteAddr"].(string); ok && remoteAddr != "" {
		fmt.Fprintf(sb, " %s", remoteAddr)
	}

	if reqBody, ok := req["body"].(string); ok && reqBody != "" {
		fmt.Fprintf(sb, "\n  %s %s", reqPrefix, formatBody(reqBody))
	}

	if userAgent, ok := req["userAgent"].(string); ok && userAgent != "" {
		fmt.Fprintf(sb, "\n  %s %s: %s", reqPrefix, userAgentField, req["userAgent"].(string))
	}

	if referrer, ok := req["referer"].(string); ok && referrer != "" {
		fmt.Fprintf(sb, "\n  %s %s: %s", reqPrefix, referrerField, referrer)
	}

	fmt.Fprintf(sb, "\n  %s %s", respPrefix, color.UnsafeGreen(log.Msg()))

	if respContentType, ok := resp["contentType"].(string); ok && respContentType != "" {
		fmt.Fprintf(sb, " %s", formatContentType(respContentType))
	}

	if size, ok := resp["size"].(float64); ok && size >= 1 {
		fmt.Fprintf(sb, " %s", color.UnsafeItalic(color.UnsafeGray(fmt.Sprintf("(%d bytes)", int(size)))))
	}

	if respBody, ok := resp["body"].(string); ok && respBody != "" {
		fmt.Fprintf(sb, "\n  %s %s", respPrefix, formatBody(respBody))
	}

	if log.Err() != "" {
		fmt.Fprintf(sb, "\n  %s %s", ErrorField, FormatError(log.Err()))
	}

	sb.WriteByte('\n')

	return sb.String()
}

func formatStatus(status float64) string {
	statusStr := color.UnsafeBold(strconv.Itoa(int(status)))

	switch {
	case status < 200:
		return color.UnsafeGray(statusStr)
	case status < 300:
		return color.UnsafeGreen(statusStr)
	case status < 400:
		return color.UnsafeBlue(statusStr)
	case status < 500:
		return color.UnsafeYellow(statusStr)
	default:
		return color.UnsafeRed(statusStr)
	}
}

func formatRequestUri(req Data) string {
	uri := ""

	if host, ok := req["host"].(string); ok {
		uri += host
	}

	if path, ok := req["path"].(string); ok {
		uri += color.UnsafeGreen(path)
	}

	return uri
}

func formatContentType(contentType string) string {
	return fmt.Sprintf("%s%s%s", leftSimpleQuote, contentType, rightSimpleQuote)
}

func formatBody(body string) string {
	body = strings.ReplaceAll(body, "\n", "\\n")

	return fmt.Sprintf("%s: %s%s%s", bodyField, bodyLeftSimpleQuote, body, bodyRightSimpleQuote)
}
