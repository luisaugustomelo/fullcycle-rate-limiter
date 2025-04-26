package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/luisaugustomelo/fullcycle-rate-limiter/internal/limiter"
	"github.com/luisaugustomelo/fullcycle-rate-limiter/internal/middleware"
)

func main() {
	_ = godotenv.Load()

	r := gin.Default()

	redisStore := limiter.NewRedisStrategy("localhost:6379")
	rl := limiter.NewRateLimiter(redisStore)

	r.Use(middleware.RateLimitMiddleware(rl))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Request accepted"})
	})

	r.Run(":8080")
}
