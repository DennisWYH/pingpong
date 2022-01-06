package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_tokens_to_pinyin(t *testing.T) {
	text := "我们都是社会主义的接班人。"

	tokens, err := Tokenizer(text)
	assert.Nil(t, err)

	var pinyins []string
	for _, val := range tokens {
		pinyin := HanziToPinyins(val)
		pinyins = append(pinyins, pinyin)
	}
	fmt.Println("tokens are", tokens, "size is", len(tokens))
	fmt.Println("pinyins are", pinyins, "size is", len(pinyins))
	assert.Equal(t, len(tokens), len(pinyins))

}

func Test_tokens_to_pinyin_with_english(t *testing.T) {
	text := "我和Mari都是社会主义的接班人。"

	tokens, err := Tokenizer(text)
	assert.Nil(t, err)

	var pinyins []string
	for _, val := range tokens {
		pinyin := HanziToPinyins(val)
		pinyins = append(pinyins, pinyin)
	}
	fmt.Println("tokens are", tokens, "size is", len(tokens))
	fmt.Println("pinyins are", pinyins, "size is", len(pinyins))
	assert.Equal(t, len(tokens), len(pinyins))

}
