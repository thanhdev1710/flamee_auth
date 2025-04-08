package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/internal/controllers"
	"github.com/thanhdev1710/flamee_auth/middlewares"
)

func AuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")

	auth.POST("/register",
		middlewares.RateLimitPerRouteAndIP(10*time.Second, 1),
		controllers.NewAuthControllers().Register,
	)

	auth.POST("/login",
		middlewares.RateLimitPerRouteAndIP(1000*time.Second, 1),
		controllers.NewAuthControllers().Login,
	)

	auth.POST("/refresh-token",
		middlewares.RateLimitPerRouteAndIP(2*time.Second, 1),
		controllers.NewAuthControllers().RefreshToken,
	)

	auth.POST("/logout", controllers.NewAuthControllers().Logout)

	auth.POST("/send-email/:email",
		middlewares.RateLimitPerRouteAndIP(20*time.Second, 1),
		controllers.NewAuthControllers().SendVerifyEmail,
	)

	auth.GET("/verify-email/:token", controllers.NewAuthControllers().VerifyEmail)

	auth.POST("/reset-password/:email",
		middlewares.RateLimitPerRouteAndIP(30*time.Second, 1),
		controllers.NewAuthControllers().SendResetPassword,
	)

	auth.POST("/change-password/:token",
		middlewares.RateLimitPerRouteAndIP(15*time.Second, 1),
		controllers.NewAuthControllers().ResetPassword,
	)

}
