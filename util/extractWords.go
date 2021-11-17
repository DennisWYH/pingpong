package util

import "fmt"

func ExtractWords(tokenizedContent []string) []string {
	// the tokenizer method output a tokenized slice with strings
	// we need to extract to tokens inside individually to build a dictionary
	var result []string
	for i, value := range tokenizedContent {
		fmt.Println("index", i)
		fmt.Println("value", value)
		if !CheckIfSymbols(value){
			result = append(result, value)
		}
	}
	return result
}