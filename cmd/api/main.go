package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/url", func(ctx *gin.Context) {
		url := ctx.DefaultQuery("url", "url")
		ctx.String(http.StatusOK, url)
	})

	r.Run()
}
