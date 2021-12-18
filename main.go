package main

import (
	"github.com/gin-gonic/gin"
	"pingpong/api"
)

type Word struct {
	Word        string `json:"word"`
	Oldword     string `json:"oldword"`
	Strokes     string `json:"strokes"`
	Pinyin      string `json:"pinyin"`
	Radicals    string `json:"radicals"`
	Explanation string `json:"explanation"`
	More        string `json:"more"`
}

func startRouting() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	// get articles
	router.GET("/articles", api.GetArticles)
	router.GET("/focusedRead", api.GetFocusedArticles)
	router.GET("/article/id/:id", api.GetArticleByID)
	router.GET("/article/grade/:grade", api.GetArticleByGrade)

	// add articles
	router.POST("/addSimpleArticle", api.AddArticle)
	router.POST("/batchAddArticles", api.BatchAddTestArticleData)

	// delete articles
	router.DELETE("/article/id/:id", api.DeleteArticleByID)
	router.DELETE("/articles", api.DeleteAllArticle)

	// update an artile
	router.PUT("/update/article/id/:id", api.UpdateArticleByID)

	// run the server
	router.Run("localhost:3456") //nolint:errcheck
}

func main() {
	startRouting()
}
