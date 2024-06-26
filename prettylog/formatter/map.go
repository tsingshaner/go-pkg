package formatter

import (
	"fmt"
	"strings"

	"github.com/tsingshaner/go-pkg/color"
)

func Map(data Data, prefix, mapKey string, parseDepth int) *strings.Builder {
	builder := &strings.Builder{}
	builder.WriteString(prefix)
	builder.WriteString(color.UnsafeBold(color.UnsafeBlue(mapKey)))

	if parseDepth > 0 {
		nestings := make(map[string]Data, len(data))

		for k, v := range data {
			if nested, ok := v.(map[string]any); ok {
				nestings[k] = nested
			} else {
				writeItem(builder, k, v)
			}
		}

		if len(nestings) > 0 {
			for key, nestRecord := range nestings {
				builder.WriteString("\n")
				builder.WriteString(Map(nestRecord,
					"  "+prefix, key, parseDepth-1).String())
			}
		}
	} else {
		for k, v := range data {
			writeItem(builder, k, v)
		}
	}

	return builder
}

func writeItem(w *strings.Builder, key string, value any) {
	w.WriteString(" ")
	w.WriteString(color.UnsafeItalic(key))
	w.WriteString(color.UnsafeDim("="))
	w.WriteString(fmt.Sprintf("%v", value))
	w.WriteString(color.UnsafeDim(","))
}
