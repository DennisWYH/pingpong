package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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

func getArticles(c *gin.Context) {
	db,_ := gorm.Open(sqlite.Open("pingpong.db"),&gorm.Config{})

	var articles []Article
	db.First(&articles)
	c.IndentedJSON(http.StatusOK, &articles)
}

func getArticleByID(c *gin.Context) {
	id := c.Param("id")
	intId,_ := strconv.Atoi(id)
	db,_ := gorm.Open(sqlite.Open("test.db"),&gorm.Config{})

	db.First(&Article{}, "ID=?", intId)
	c.IndentedJSON(http.StatusOK, &Article{})
}

// API: curl -X POST localhost:3456/addArticle/hello/world
func addArticle(c *gin.Context){
	fmt.Println("served by addArticle handler.")
	var newArticle Article
	newArticle.Title = c.Param("title")
	newArticle.Content = c.Param("content")

	// Add the new article to the db table.
	db,_ := gorm.Open(sqlite.Open("pingpong.db"), &gorm.Config{})
	db.Create(&newArticle)

	// show the albums table after adding an entry
	var articles []Article
	db.Find(&articles)
	c.IndentedJSON(http.StatusCreated, &articles)

}