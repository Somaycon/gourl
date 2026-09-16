package main

import (
	"os"

	"github.com/Somaycon/gourl.git/internal/handler"
	"github.com/Somaycon/gourl.git/internal/model"
	"github.com/Somaycon/gourl.git/internal/repository"
	"github.com/Somaycon/gourl.git/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	dsn := os.Getenv("DSN")

	db, _ := gorm.Open(postgres.Open(dsn))
	db.AutoMigrate(&model.Url{})

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	urlRepoisoty := repository.NewUrlRepository(db, rdb)
	urlService := service.NewUrlService(urlRepoisoty, os.Getenv("BASE_URL"))
	urlHandler := handler.NewUrlHandler(urlService)

	r := gin.Default()
	r.POST("/url", urlHandler.Create)
	r.GET("/:short", urlHandler.Redirect)

	r.Run()

}
