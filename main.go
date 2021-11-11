package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"pingpong/db"
)
// lookupC2C given a chinese word
// look it up in word.json and return the explanation
func lookupC2C(c string) {
	dictionaryPath := "./chinese-xinhua/data/word.json"
	file, _ := ioutil.ReadFile(dictionaryPath)
	var data = make([]Word, 100)
	err := json.Unmarshal(file, &data)
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < len(data); i++ {
		if data[i].Word == c {
			fmt.Println(data[i].Explanation)
		}
	}

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
	router := gin.Default()
	router.GET("/articles", getArticles)
	router.GET("/articles/:id", getArticleByID)

	router.POST("/addArticle/:title/:content", addArticle)
	router.Run("localhost:3456")
}
