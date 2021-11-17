package util

import (
	"github.com/jcramb/cedict"
)

func Cn_en_lookup(hanzi string) []string {
	d := cedict.New()
	entry := d.GetByHanzi(hanzi)
	enMeaning := entry.Meanings
	return enMeaning
}