package main

import (
	"github.com/gin-gonic/gin"
)

type Word struct {
	Word        string `json:"word"`
	Oldword     string `json:"oldword"`
	Strokes     string `json:"strokes"`
	Pinyin      string `json:"pinyin"`
	Radicals    string `json:"radicals"`
	Explanation string `json:"explanation"'`
	More        string `json:"more"`
}

func startRouting(){
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/articles", getArticles)
	router.GET("/article/id/:id", getArticleByID)
	router.GET("/article/grade/:grade", getArticleByGrade)
	router.POST("/addArticle", addArticle)
	router.Run("localhost:3456")
}


func main() {
	//CreateDBTables()
//	cn_en.Cn_en_lookup("总共")
	startRouting()
}
