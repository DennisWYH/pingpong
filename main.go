package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func runServer() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	router.Run(":8083")
}

type Word struct {
	Word        string `json:"word"`
	Oldword     string `json:"oldword"`
	Strokes     string `json:"strokes"`
	Pinyin      string `json:"pinyin"`
	Radicals    string `json:"radicals"`
	Explanation string `json:"explanation"'`
	More        string `json:"more"`
}

func main() {
	dictionaryPath := "./chinese-xinhua/data/word.json"
	file, _ := ioutil.ReadFile(dictionaryPath)
	var data = make([]Word, 100)
	err := json.Unmarshal(file, &data)
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < len(data); i++ {
		if data[i].Word == "你" {
			fmt.Println(data[i].Explanation)
		}
	}
}
