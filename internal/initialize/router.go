package initialize

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/internal/routes"
	"github.com/thanhdev1710/flamee_auth/middlewares"
)

func InitRouter() *gin.Engine {
	var r *gin.Engine
	if global.Config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	} else {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		r = gin.Default()
	}

	// 1. Chặn bot / request đáng ngờ (càng sớm càng tốt)
	r.Use(middlewares.BlockSuspiciousUserAgents())

	// 2. Logging để ghi nhận request sau khi đã qua bộ lọc ban đầu
	r.Use(middlewares.RequestLogger())

	// 3. Thiết lập các header bảo mật
	middlewares.Helmet(r)

	// 4. CORS - nên đặt trước mọi route logic
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{global.Url.UrlFrontEnd}, // Cập nhật đúng URL frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-API-KEY"},
		ExposeHeaders:    []string{"X-Total-Count"}, // Để expose thêm header cho frontend
		AllowCredentials: true,                      // Quan trọng: Cho phép gửi cookies
	}))

	// 5. Kiểm tra API key
	r.Use(middlewares.CheckAPIKey())

	// 6. Tin tưởng proxy (nếu có dùng reverse proxy như Nginx)
	r.SetTrustedProxies(nil)

	// 7. Định tuyến
	adminRouter := routes.RouterGroupApp.Admin
	userRouter := routes.RouterGroupApp.User
	authRouter := routes.RouterGroupApp.Auth

	MainGroup := r.Group("/api/v1")
	{
		MainGroup.GET("/check-status")
	}
	{
		authRouter.InitAuthRouter(MainGroup)
	}
	{
		adminRouter.Post.InitPostRouter(MainGroup)
	}
	{
		userRouter.Post.InitInteractionRouter(MainGroup)
		userRouter.Post.InitPostRouter(MainGroup)
		userRouter.Profile.InitProfileRouter(MainGroup)
		userRouter.Profile.InitFollowRouter(MainGroup)
		userRouter.Search.InitSearchRouter(MainGroup)
	}

	return r
}
