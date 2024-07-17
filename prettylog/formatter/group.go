package formatter

import (
	"strings"

	"github.com/tsingshaner/go-pkg/color"
)

func Groups(groups []Group, prefix, mapKey string) *strings.Builder {
	sb := &strings.Builder{}
	sb.WriteString(prefix)
	sb.WriteString(color.UnsafeBold(color.UnsafeBlue(mapKey)))

	for i, group := range groups {
		sb.WriteString(color.UnsafeBold(color.UnsafeBlue(group.Key)))

		if len(group.Value) == 0 {
			if i < len(groups)-1 {
				sb.WriteString(color.UnsafeBold(color.UnsafeRed(".")))
			}

			continue
		}

		for k, v := range group.Value {
			writeItem(sb, k, v)
		}

		if i < len(groups)-1 {
			sb.WriteByte('\n')
			prefix = "  " + prefix
			sb.WriteString(prefix)
		}
	}

	return sb
}
