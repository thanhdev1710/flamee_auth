package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/global"
	"go.uber.org/zap"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Tiếp tục request
		c.Next()

		// Sau khi xử lý xong
		duration := time.Since(start)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		userId := c.GetString("userId")

		global.Logger.Info("HTTP Request",
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("clientIP", clientIP),
			zap.String("userAgent", userAgent),
			zap.String("userId", userId),
			zap.Duration("latency", duration),
		)
	}
}
