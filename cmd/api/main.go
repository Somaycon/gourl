package main

import (
	"os"
	"time"

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
	_ = godotenv.Load()
	dsn := os.Getenv("DSN")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&model.Url{}); err != nil {
		panic("Falha ao executar AutoMigrate: " + err.Error())
	}

	postgresDB, err := db.DB()
	if err != nil {
		panic("Falha ao obter *sql.DB: " + err.Error())
	}

	postgresDB.SetMaxIdleConns(10)
	postgresDB.SetMaxOpenConns(100)
	postgresDB.SetConnMaxLifetime(5 * time.Minute)

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
