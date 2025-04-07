package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/global"
)

func CheckAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy API Key từ header
		apiKey := c.GetHeader("X-API-Key")

		// Kiểm tra API Key hợp lệ (sử dụng biến môi trường hoặc key hardcoded)
		if apiKey != global.Config.ApiKey { // Hoặc một giá trị cố định
			c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden"})
			c.Abort()
			return
		}

		// Nếu API Key hợp lệ, tiếp tục xử lý
		c.Next()
	}
}
