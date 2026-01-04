package formatting

import (
	"strings"
	"unicode/utf8"
)

func PadString(s string, totalLen int) string {
	sLen := utf8.RuneCountInString(s)

	if sLen >= totalLen {
		return s
	}

	padCount := totalLen - sLen

	padding := strings.Repeat(" ", padCount/2)

	postfix := ""

	if padCount%2 == 1 {
		postfix = " "
	}

	return padding + s + padding + postfix
}
