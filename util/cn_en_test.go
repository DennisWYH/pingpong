package util

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_cn_en(t *testing.T) {
	word1 := "中国"
	enLookup1, err := Cn_en_lookup(word1)
	assert.Nil(t, err)
	assert.Equal(t, []string{"China"}, enLookup1)

	word2 := "说话"
	enLookup2, err := Cn_en_lookup(word2)
	assert.Nil(t, err)
	assert.Equal(t, []string{"to speak", "to say", "to talk", "to gossip", "to tell stories", "talk", "word"}, enLookup2)

	// words that are potentially not in the dictionary ()
	// current if you look up a word that is not in the dictionary it panics.
	word3 := "已大"
	enLookup3, err := Cn_en_lookup(word3)
	assert.Equal(t, []string(nil), enLookup3)
	assert.Equal(t, errors.New("no result lookup for this word"), err)
}
