package util

import (
	"errors"
	"github.com/jcramb/cedict"
)

func Cn_en_lookup(hanzi string) ([]string, error) {
	d := cedict.New()
	entry := d.GetByHanzi(hanzi)
	// if no meaning is found, return nil
	if entry == nil {
		return nil, errors.New("no result lookup for this word")
	}
	enMeaning := entry.Meanings
	return enMeaning, nil
}
