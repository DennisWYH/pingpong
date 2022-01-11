package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"pingpong/api"
	"strconv"
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

func castInt(val string) (int, error) {
	return strconv.Atoi(val)
}
func getValueAtIndex(val []string, index int) string {
	return val[index]
}
func getFirst160(val string) string {
	return val[:160]
}
func lenBiggerThan160(val string) bool {
	if len(val) > 160 {
		return true
	} else {
		return false
	}
}

func startRouting() {
	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"castInt":          castInt,
		"getValueAtIndex":  getValueAtIndex,
		"getFirst160":      getFirst160,
		"lenBiggerThan160": lenBiggerThan160,
	})
	router.LoadHTMLGlob("static/templates/*")

	// design previews
	router.GET("/design1", api.GetDesign1Handler)

	// get articles
	router.GET("/articles", api.GetArticlesHandler)
	router.GET("/focusedRead/id/:articleID", api.GetFocusedArticlesHandler)
	router.GET("/article/id/:articleID", api.GetArticleByIDHandler)
	router.GET("/article/grade/:grade", api.GetArticleByGradeHandler)

	// add articles
	router.POST("/addSimpleArticle", api.AddArticleHandler)

	// delete articles
	router.DELETE("/article/id/:articleID", api.DeleteArticleByIDHandler)
	router.DELETE("/articles", api.DeleteAllArticleHandler)

	// update an article
	router.PUT("/update/article/id/:articleID", api.UpdateArticleByIDHandler)

	// serve static files
	router.Static("./static", "./static")
	// run the server
	router.Run("localhost:3456") //nolint:errcheck
}

func main() {
	startRouting()
	//database.CreateDBTables()
}
