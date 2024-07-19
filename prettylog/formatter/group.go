package formatter

import (
	"strings"

	"github.com/tsingshaner/go-pkg/color"
	"github.com/tsingshaner/go-pkg/util/slices"
)

func Groups(groups []Group, prefix, mapKey string) *strings.Builder {
	sb := &strings.Builder{}
	groups = cleanGroup(groups)
	if len(groups) == 0 {
		return sb
	}

	sb.WriteByte('\n')
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

func cleanGroup(groups []Group) []Group {
	lastEmptyGroupIndex := slices.LastIndexFunc(groups, func(g Group) bool {
		return len(g.Value) > 0
	})

	if lastEmptyGroupIndex == -1 {
		return nil
	}

	return groups[:lastEmptyGroupIndex+1]
}
