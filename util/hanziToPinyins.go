package util

import (
	"github.com/Lofanmi/pinyin-golang/pinyin"
)

func HanziToPinyins(hanzi string) string {
	dict := pinyin.NewDict()
	// wo3, he2 shi2 neng2 bao4 fu4?
	result := dict.Sentence(hanzi).Unicode()
	return result
}
