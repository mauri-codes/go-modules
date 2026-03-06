package utils

import (
	"strings"
	"unicode"
)

func CamelToKebab(s string) string {
	var result strings.Builder

	for index, r := range s {
		if unicode.IsUpper(r) {
			if index > 0 {
				result.WriteRune('-')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}
