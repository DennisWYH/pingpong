package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mozillazg/go-pinyin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)
type Article struct {
	Title         string
	Content       string
	//Tags          []string
	//WordCount     int64
	//Grade         string
	//NumberOfRead  int64
	//NumberOfFlash int64
	gorm.Model
}

// API: localhost:3456/articles
func getArticles(c *gin.Context) {
	db,_ := gorm.Open(sqlite.Open("pingpong.db"),&gorm.Config{})

	var articles []Article
	db.First(&articles)

	c.HTML(http.StatusOK, "viewArticles.tmpl", gin.H{
		"articles": &articles,
	})
}

// API: localhost:3456/article/:ID
func getArticleByID(c *gin.Context) {
	id := c.Param("id")
	intId,_ := strconv.Atoi(id)
	db,_ := gorm.Open(sqlite.Open("pingpong.db"),&gorm.Config{})

	var article *Article
	db.First(&article, "ID=?", intId)

	articleStruct := *article
	content := articleStruct.Content
	a := pinyin.NewArgs()
	a.Style = pinyin.Tone
	contentPinyins := pinyin.Pinyin(content, a)
	fmt.Println("pin yin is ", contentPinyins)

	slicedContent := strings.Split(content, "")

	hanziPinyins := make(map[string][]string)
	for i:=0; i< len(slicedContent);i++ {
		key := slicedContent[i]
		value := contentPinyins[i]
		hanziPinyins[key] = value
	}

	c.HTML(http.StatusOK, "viewArticleById.tmpl", gin.H{
		"hanzi": content,
		"hanziPinyins" : hanziPinyins,
	})
}

// API: curl -X POST -H "Content-Type: application/x-www-form-urlencoded"
//  -d "title=new&content=entry" localhost:3456/addArticle
// gin context documentation: https://pkg.go.dev/github.com/gin-gonic/gin#section-readme
func addArticle(c *gin.Context){
	fmt.Println("served by addArticle handler.")
	var newArticle Article

	newArticle.Title = c.PostForm("title")
	//Todo: the content from user input has to be chinese,
	// for later pinyin convert.
	//Todo: what if there are English words in the paragraph...
	// solution: create a map, when English reconized, put blank in there or display english itself.
	newArticle.Content = c.PostForm("content")

	// Add the new article to the db table.
	db,_ := gorm.Open(sqlite.Open("pingpong.db"), &gorm.Config{})
	db.Create(&newArticle)

	// show the albums table after adding an entry
	var articles []Article
	db.Find(&articles)

	c.HTML(http.StatusCreated, "viewArticles.tmpl", gin.H{
		"articles": &articles,
	})
}