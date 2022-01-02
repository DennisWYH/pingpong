package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_cn_en(t *testing.T) {
	word1 := "中国"
	enLookup1 := Cn_en_lookup(word1)
	assert.Equal(t, []string{"China"}, enLookup1)

	word2 := "说话"
	enLookup2 := Cn_en_lookup(word2)
	assert.Equal(t, []string{"to speak", "to say", "to talk", "to gossip", "to tell stories", "talk", "word"}, enLookup2)

	// words that are potentially not in the dictionary ()
	// current if you look up a word that is not in the dictionary it panics.
	word3 := "已大"
	enLookup3 := Cn_en_lookup(word3)
	assert.Equal(t, "xxx", enLookup3)

}
