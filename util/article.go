package util

import "strings"

func CountWords(content string) int64 {
	return int64(len(strings.Fields(content)))
}
