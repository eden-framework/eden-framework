package builder

import (
	"strings"
)

func HolderRepeat(length int) string {
	if length > 1 {
		return strings.Repeat("?,", length-1) + "?"
	}
	return "?"
}

func quote(n string) string {
	return "`" + n + "`"
}
