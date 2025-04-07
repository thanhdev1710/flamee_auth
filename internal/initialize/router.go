package initialize

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/internal/routes"
	"github.com/thanhdev1710/flamee_auth/middlewares"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	middlewares.Helmet(r)
	// Thêm middleware CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://yourtrustedwebsite.com"},          // Cho phép truy cập chỉ từ domain đáng tin cậy
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},            // Xác định các phương thức HTTP được phép
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Cho phép các headers cần thiết cho CORS
		ExposeHeaders:    []string{"X-Total-Count"},                           // Cho phép client truy cập một số headers đặc biệt
		AllowCredentials: true,                                                // Cho phép gửi cookies và thông tin xác thực
		MaxAge:           12 * time.Hour,                                      // Cấu hình thời gian cache preflight request
	}))

	// Thêm middleware kiểm tra API key
	r.Use(middlewares.CheckAPIKey())

	// Đăng ký các route
	routes.AuthRoutes(r)
	routes.ProtectRoutes(r)

	return r
}
