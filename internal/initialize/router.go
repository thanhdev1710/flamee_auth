package initialize

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/internal/routes"
	"github.com/thanhdev1710/flamee_auth/middlewares"
)

func InitRouter() *gin.Engine {
	if global.Config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// 1. Chặn bot / request đáng ngờ (càng sớm càng tốt)
	r.Use(middlewares.BlockSuspiciousUserAgents())

	// 2. Logging để ghi nhận request sau khi đã qua bộ lọc ban đầu
	r.Use(middlewares.RequestLogger())

	// 3. Thiết lập các header bảo mật
	middlewares.Helmet(r)

	// 4. CORS - nên đặt trước mọi route logic
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://yourtrustedwebsite.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"X-Total-Count"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 5. Kiểm tra API key
	r.Use(middlewares.CheckAPIKey())

	// 6. Tin tưởng proxy (nếu có dùng reverse proxy như Nginx)
	r.SetTrustedProxies(nil)

	// 7. Định tuyến
	routes.AuthRoutes(r)
	routes.ProtectRoutes(r)

	return r
}
