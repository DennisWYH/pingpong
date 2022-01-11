package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"pingpong/util"
	"strconv"
)

// GetFocusedArticlesHandler handles the request and renders viewFocusedRead html
// API: curl localhost:3456/focusedRead/id/:articleID
func GetFocusedArticlesHandler(c *gin.Context) {
	db, _ := gorm.Open(sqlite.Open("pingpong.db"), &gorm.Config{})
	id := c.Param("articleID")
	intID, _ := strconv.Atoi(id)

	var article *Article
	db.First(&article, "ID=?", intID)

	var lookups []Lookup
	db.Where("article_id = ?", article.ID).Find(&lookups)

	// for each article content, we first tokenize it
	tokens, err := util.Tokenizer(article.Content)
	pinyins := Tokens_to_pinyins(tokens)

	if err != nil {
		fmt.Print("There is an error in tokenizing the article content", err)
	}

	c.HTML(http.StatusOK, "design1.html", gin.H{
		"tokens":  tokens,
		"pinyins": pinyins,
		"article": &article,
		"lookups": &lookups,
	})
}
