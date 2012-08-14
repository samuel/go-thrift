package thrift

import (
	"strings"
	"unicode"
)

func camelCase(s string) string {
	prev := '_'
	return strings.Map(
		func(r rune) rune {
			if r == '_' {
				prev = r
				return -1
			}
			if prev == '_' {
				prev = r
				return unicode.ToUpper(r)
			}
			prev = r
			return r
		}, s)
}
