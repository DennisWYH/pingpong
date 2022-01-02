package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Tokenizer(t *testing.T) {
	text1 := "我是一个中国人"
	result1, err := Tokenizer(text1)
	if err != nil {
		t.Errorf("Error in tokenizing text: %s", text1)
	}
	assert.Equal(t, []string{"我", "是", "一个", "中国", "人"}, result1)

	text2 := "明天天气会更好"
	result2, err := Tokenizer(text2)
	if err != nil {
		t.Errorf("Error in tokenizing text: %s", text2)
	}
	assert.Equal(t, []string{"明天", "天气", "会", "更好"}, result2)

	// special symbols in the sentence can be tokenized correctly
	text3 := "你们说：怎么安排？"
	result3, err := Tokenizer(text3)
	if err != nil {
		t.Errorf("Error in tokenizing text: %s", text3)
	}
	assert.Equal(t, []string{"你们", "说", "：", "怎么", "安排", "？"}, result3)

	// English word single appearance in Chinese sentence can be tokenized.
	text4 := "我们明天一起camping你觉得怎么样？"
	result4, err := Tokenizer(text4)
	if err != nil {
		t.Errorf("Error in tokenizing text: %s", text4)
	}
	assert.Equal(t, []string{"我们", "明天", "一起", "camping", "你", "觉得", "怎么样", "？"}, result4)

	// English words multiple appearances in Chinese sentence can be tokenized.
	text5 := "他说：'How are you?'"
	result5, err := Tokenizer(text5)
	if err != nil {
		t.Errorf("Error in tokenizing text: %s", text5)
	}
	assert.Equal(t, []string{"他", "说", "：", "'", "How", " ", "are", " ", "you", "?", "'"}, result5)
}
