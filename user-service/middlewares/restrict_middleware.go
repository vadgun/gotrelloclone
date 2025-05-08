package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimitMiddleware() gin.HandlerFunc {
	limiter := rate.NewLimiter(1, 1) // 1 solicitud por segundo con un burst de 1
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Demasiadas peticiones!"})
			c.Abort()
			return
		}
		c.Next()
	}
}
