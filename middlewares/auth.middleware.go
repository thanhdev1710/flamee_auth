package middlewares

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/internal/models"
	"github.com/thanhdev1710/flamee_auth/pkg/utils"
	"gorm.io/gorm"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Vui lòng đăng nhập",
			})
			c.Abort()
			return
		}

		tokenStr := parts[1]

		// Kiểm tra token hợp lệ
		claims, err := utils.ValidateToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("userId", claims.Subject)
		c.Set("role", claims.Role)
		c.Set("jwt", authHeader)

		c.Next()
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
				"status":  "error",
				"message": "Bạn không có quyền truy cập",
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

func VerifyAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := utils.GetUserId(c)
		// Khai báo biến user để lưu dữ liệu người dùng
		var user models.User

		// Tìm người dùng theo userId
		err := global.Pdb.Where("id = ?", userId).First(&user).Error
		if err != nil {
			// Nếu lỗi là không tìm thấy người dùng
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Không tìm thấy tài khoản"})
			} else {
				// Nếu có lỗi khác khi truy vấn DB
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Lỗi server"})
			}
			c.Abort()
			return
		}

		// Kiểm tra xem người dùng đã xác thực email chưa
		if !user.IsVerified || user.Status != global.User.Active {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Tài khoản này chưa xác thực hoặc đã xoá"})
			c.Abort()
			return
		}

		// Tiến hành các handler tiếp theo nếu email đã được xác thực
		c.Next()
	}
}
