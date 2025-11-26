package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/pkg/utils"
	"go.uber.org/zap"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {

		// ----- SKIP LOGS ROUTE -----
		skip := map[string]bool{
			"/api/v1/logs": true,
		}
		if skip[c.Request.URL.Path] {
			c.Next()
			return
		}

		start := time.Now()

		// tiếp tục xử lý request
		c.Next()

		// sau khi xử lý xong
		duration := time.Since(start)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		userId := utils.GetUserId(c)

		// -------------------------
		// Chọn level theo status
		// -------------------------
		var logFn func(msg string, fields ...zap.Field)
		var levelLabel string

		switch {
		case statusCode >= 500:
			logFn = global.Logger.Error
			levelLabel = "ERROR"
		case statusCode >= 400:
			logFn = global.Logger.Warn
			levelLabel = "WARN"
		default:
			logFn = global.Logger.Info
			levelLabel = "INFO"
		}

		// Log ra
		logFn("HTTP Request",
			zap.String("level", levelLabel),
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
