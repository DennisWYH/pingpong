package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mozillazg/go-pinyin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"pingpong/util"
	"strconv"
	"strings"
)

type Article struct {
	Title         string
	Content       string
	Grade         string
	//Tags          []string
	//WordCount     int64
	//NumberOfRead  int64
	//NumberOfFlash int64
	gorm.Model
}

// API: localhost:3456/articles
func getArticles(c *gin.Context) {
	db,_ := gorm.Open(sqlite.Open("pingpong.db"),&gorm.Config{})

	var articles []Article
	db.Find(&articles)

	c.IndentedJSON(http.StatusOK, &articles)
	c.HTML(http.StatusOK, "viewArticles.tmpl", gin.H{
		"articles": &articles,
	})
}

// API: curl -X DELETE localhost:3456/article/id/:id
func deleteArticleByID(c *gin.Context) {
	db,_ := gorm.Open(sqlite.Open("pingpong.db"),&gorm.Config{})

	var article *Article
	id := c.Param("id")
	intId,_ := strconv.Atoi(id)
	db.Delete(&article, intId)

	var articles *[]Article
	db.Find(&articles)
	c.IndentedJSON(http.StatusOK, &articles)
	c.HTML(http.StatusOK, "viewArticles.tmpl", gin.H{
		"articles": &articles,
	})
}

// API: curl -X PUT -d "content=?" localhost:3456/update/article/id/:id
func updateArticleByID(c *gin.Context) {
	db,_ := gorm.Open(sqlite.Open("pingpong.db"),&gorm.Config{})

	var article *Article
	id := c.Param("id")
	intId,_ := strconv.Atoi(id)

	db.First(&article, "ID=?", intId)

	article.Content = c.PostForm("content")
	db.Save(&article)

	c.IndentedJSON(http.StatusOK, &article)
	c.HTML(http.StatusOK, "viewArticles.tmpl", gin.H{
		"articles": &article,
	})
}

// API: localhost:3456/article/id/:id
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

// API: localhost:3456/article/grade/:grade
func getArticleByGrade(c *gin.Context) {
	grade := c.Param("grade")
	db,_ := gorm.Open(sqlite.Open("pingpong.db"),&gorm.Config{})

	var articles *[]Article
	db.Find(&articles, "Grade=?", grade)


	var hanzis []string
	var pinyins []string
	var tokenizedContents [][]string
	var words []string
	var wordsEns [][]string

	for _, article := range *articles {
		articleStruct := article
		content := articleStruct.Content

		hanzis = append(hanzis, content)
		pinyins = append(pinyins, util.HanziToPinyins(content))
		tokenizedContent := util.Tokenizer(content)
		tokenizedContents = append(tokenizedContents, tokenizedContent)
		words = util.ExtractWords(tokenizedContent)

	}
	for _, word := range words {
		wordsEn := util.Cn_en_lookup(word)
		wordsEns = append(wordsEns, wordsEn)
	}

	c.HTML(http.StatusOK, "viewArticleByGrade.tmpl", gin.H{
		"hanzis": hanzis,
		"pinyins": pinyins,
		"tokenizedContents" : tokenizedContents,
		"words": words,
		"wordsEns": wordsEns,
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
	newArticle.Grade = c.PostForm("grade")

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