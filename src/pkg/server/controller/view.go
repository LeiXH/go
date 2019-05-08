package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GinIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "html/index.html", gin.H{})
}

func GinList(c *gin.Context) {
	c.HTML(http.StatusOK, "html/list.html", gin.H{})
}
