package middleware

import (
	"net/http"
	"strings"

	"github.com/luisaugustomelo/fullcycle-rate-limiter/internal/limiter"

	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware(l *limiter.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		token := strings.TrimSpace(c.GetHeader("API_KEY"))
		key := l.GetKey(ip, token)
		limit := l.GetLimit(token)

		allowed, err := l.Store.AllowRequest(key, limit, 1)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
			c.Abort()
			return
		}

		if !allowed {
			l.Store.SetBlock(key, l.BlockTime)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message": "you have reached the maximum number of requests or actions allowed within a certain time frame",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
