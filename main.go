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
	router.LoadHTMLGlob("static/templates/*")

	// get articles
	router.GET("/articles", api.GetArticlesHandler)
	router.GET("/focusedRead", api.GetFocusedArticlesHandler)
	router.GET("/article/id/:id", api.GetArticleByIDHandler)
	router.GET("/article/grade/:grade", api.GetArticleByGradeHandler)

	// add articles
	router.POST("/addSimpleArticle", api.AddArticleHandler)
	router.POST("/batchAddArticles", api.BatchAddTestArticleDataHandler)

	// delete articles
	router.DELETE("/article/id/:id", api.DeleteArticleByIDHandler)
	router.DELETE("/articles", api.DeleteAllArticleHandler)

	// update an artile
	router.PUT("/update/article/id/:id", api.UpdateArticleByIDHandler)

	// serve static files
	router.Static("./static", "./static")
	// run the server
	router.Run("localhost:3456") //nolint:errcheck
}

func main() {
	startRouting()
	//api.GetLookups()
}
