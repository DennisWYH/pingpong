package util

import "github.com/xujiajun/gotokenizer"

// Tokenizer takes in a string of chinese and tokenize it into slice of strings
func Tokenizer(text string) ([]string, error) {
	dictPath := "/Users/yunhaiwang/go/pkg/mod/github.com/xujiajun/gotokenizer@v1.1.0/data/zh/dict.txt"
	// NewMaxMatch default wordFilter is NumAndLetterWordFilter
	mm := gotokenizer.NewMaxMatch(dictPath)
	// load dict
	err := mm.LoadDict()
	if err != nil {
		return nil, err
	}

	tokenizedStrings, _ := mm.Get(text)
	return tokenizedStrings, nil
}
