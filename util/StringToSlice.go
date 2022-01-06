package util

import "strings"

func StringToSlice(text string) []string {
	return strings.Split(text, ",")
}
