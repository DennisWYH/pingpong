package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Tokenizer(t *testing.T) {
	text := "我是一个中国人"
	result, err := Tokenizer(text)
	if err != nil {
		t.Errorf("Error in tokenizing text: %s", text)
	}
	assert.Equal(t, []string{"我", "是", "一个", "中国", "人"}, result)

	text = "明天天气会更好"
	result, err = Tokenizer(text)
	if err != nil {
		t.Errorf("Error in tokenizing text: %s", text)
	}
	assert.Equal(t, []string{"明天", "天气", "会", "更好"}, result)

	// special symbols in the sentence can be tokenized correctly
	text = "你们说：怎么安排？"
	result, err = Tokenizer(text)
	if err != nil {
		t.Errorf("Error in tokenizing text: %s", text)
	}
	assert.Equal(t, []string{"你们", "说", "：", "怎么", "安排", "？"}, result)

	// English word single appearance in Chinese sentence can be tokenized.
	text = "我们明天一起camping你觉得怎么样？"
	result, err = Tokenizer(text)
	if err != nil {
		t.Errorf("Error in tokenizing text: %s", text)
	}
	assert.Equal(t, []string{"我们", "明天", "一起", "camping", "你", "觉得", "怎么样", "？"}, result)

	// English words multiple appearances in Chinese sentence can be tokenized.
	text = "他说：'How are you?'"
	result, err = Tokenizer(text)
	if err != nil {
		t.Errorf("Error in tokenizing text: %s", text)
	}
	assert.Equal(t, []string{"他", "说", "：", "'", "How", " ", "are", " ", "you", "?", "'"}, result)
}
