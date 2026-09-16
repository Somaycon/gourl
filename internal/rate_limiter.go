package internal

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimiter(rdb *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientIP := ctx.ClientIP()
		key := "rate_limit" + clientIP
		c := ctx.Request.Context()

		requests, err := rdb.Incr(c, key).Result()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			ctx.Abort()
			return
		}

		if requests == 1 {
			rdb.Expire(c, key, window)
		}

		if requests > int64(limit) {
			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too Many Requests",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
