package util

import (
	"unicode"
)

func CheckIfChineseChar(char rune) bool {
	if unicode.Is(unicode.Han, char) {
		return true
	}
	return false
}
