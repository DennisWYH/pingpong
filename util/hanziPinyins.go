package util

import (
	"github.com/mozillazg/go-pinyin"
	"strings"
)

func MakeHanziWithPinyins (hanzis string) map[string][]string {
	a := pinyin.NewArgs()
	a.Style = pinyin.Tone
	contentPinyins := pinyin.Pinyin(hanzis, a)

	slicedContent := strings.Split(hanzis, "")

	hanziPinyins := make(map[string][]string)
	for i:=0; i< len(slicedContent);i++ {
		key := slicedContent[i]
		value := contentPinyins[i]
		hanziPinyins[key] = value
	}
	return hanziPinyins
}
