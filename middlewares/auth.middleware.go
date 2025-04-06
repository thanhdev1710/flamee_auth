package middlewares

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/pkg/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Lấy token từ header Authorization
		tokenStr := strings.Split(ctx.GetHeader("Authorization"), " ")[1]

		// Kiểm tra token hợp lệ
		claims, err := utils.ValidateToken(tokenStr)
		if err != nil {
			// Nếu token không hợp lệ, trả về lỗi
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			ctx.Abort()
			return
		}

		ctx.Set("userId", claims.Subject)
		ctx.Set("role", claims.Role)

		ctx.Next()
	}
}

func RoleRequired(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy thông tin role từ context (có thể đã được set ở AuthMiddleware)
		role := c.GetString("role")

		// Kiểm tra xem role có nằm trong danh sách roles được phép không
		if !contains(roles, role) {
			// Nếu không có quyền truy cập, trả về lỗi 403 (Forbidden)
			c.JSON(http.StatusForbidden, gin.H{
				"message": "You do not have permission to access this resource",
			})
			c.Abort() // Ngừng tiếp tục xử lý request
			return
		}

		// Nếu role hợp lệ, tiếp tục xử lý request
		c.Next()
	}
}

// Hàm kiểm tra xem role có nằm trong danh sách roles không
func contains(roles []string, role string) bool {
	return slices.Contains(roles, role)
}
