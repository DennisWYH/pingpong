package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

func GetArticlesHandler(c *gin.Context) {
	db, _ := gorm.Open(sqlite.Open("pingpong.db"), &gorm.Config{})

	var articles []Article
	db.Find(&articles)

	c.HTML(http.StatusOK, "articles.html", gin.H{
		"articles": articles,
	})
}
