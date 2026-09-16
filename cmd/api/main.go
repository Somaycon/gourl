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
	gorm.Model
	Code     string `json:"code" gorm:"uniqueIndex"`
	ShortUrl string `json:"short_url" binding:"required, url"`
	Url      string `json:"url"`
	Clicks   int    `gorm:"default:0" json:"clicks"`
}

type UrlRequest struct {
	Url string `json:"url"`
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
		var request UrlRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		var code string

		for {
			generated, err := GenerateCode(6)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "Error to generate code",
				})
			}
			var count int64
			db.Model(&Url{}).Where("code = ?", generated).Count(&count)
			if count == 0 {
				code = generated
				break
			}
		}
		shortUrl := os.Getenv("BASE_URL") + code

		newUrl := Url{
			Code:     code,
			ShortUrl: shortUrl,
			Url:      request.Url,
		}
		db.Create(&newUrl)
		ctx.JSON(http.StatusOK, newUrl.ShortUrl)
	})

	r.GET("/:short", func(ctx *gin.Context) {
		short := ctx.Param("short")
		var url Url
		if err := db.Where("code = ?", short).First(&url).Error; err != nil {
			ctx.String(http.StatusNotFound, "url not found")
			return
		}
		db.Model(&url).UpdateColumn("clicks", gorm.Expr("clicks + ?", 1))
		ctx.Redirect(http.StatusFound, url.Url)
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
