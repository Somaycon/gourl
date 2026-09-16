package main

import (
	"crypto/rand"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Url struct {
	Code string `json:"code"`
	Url  string `json:"url"`
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func main() {
	r := gin.Default()

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/url", func(ctx *gin.Context) {
		url := ctx.DefaultQuery("url", "url")
		code, err := GenerateCode(6)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		newUrl := Url{
			Code: code,
			Url:  url,
		}
		ctx.JSON(http.StatusOK, newUrl)
	})

	r.Run()
}

func GenerateCode(length int) (string, error) {
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[num.Int64()]
	}
	return string(b), nil
}
