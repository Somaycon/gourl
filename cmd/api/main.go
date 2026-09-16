package main

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Url struct {
	*gorm.Model
	Code     string `json:"code"`
	ShortUrl string `json:"short_url"`
	Url      string `json:"url"`
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	dsn := os.Getenv("DSN")

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Url{})

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
		if err := db.Where("code = ?", code).First(&Url{}).Error; err == nil {
			code, err = GenerateCode(6)
			if err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}
		}
		shortUrl := os.Getenv("BASE_URL") + code

		newUrl := Url{
			Code:     code,
			ShortUrl: shortUrl,
			Url:      url,
		}
		db.Create(&newUrl)
		ctx.JSON(http.StatusOK, newUrl.ShortUrl)
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
