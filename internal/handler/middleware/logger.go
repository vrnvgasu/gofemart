package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/vrnvgasu/gofemart/internal/logger"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		logger.Log.Infow("http request",
			"uri", c.Request.RequestURI,
			"method", c.Request.Method,
			"duration", duration,
			"status", c.Writer.Status(),
			"size", c.Writer.Size(),
		)
	}
}
