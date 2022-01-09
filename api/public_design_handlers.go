package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetDesign1Handler(c *gin.Context) {
	c.HTML(http.StatusOK, "design1.html", gin.H{"": ""})
}
