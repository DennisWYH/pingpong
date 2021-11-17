package util

import (
	"fmt"
	"github.com/jcramb/cedict"
)

func Cn_en_lookup(hanzi string) []string {
	d := cedict.New()
	entry := d.GetByHanzi(hanzi)
	enMeaning := entry.Meanings
	fmt.Println(enMeaning)
	return enMeaning
}