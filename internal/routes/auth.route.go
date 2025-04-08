package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/internal/controllers"
	"github.com/thanhdev1710/flamee_auth/middlewares"
)

func AuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	authController := controllers.NewAuthControllers()

	auth.POST("/register",
		middlewares.RateLimitPerRouteAndIP(1, 30*time.Second, 3),
		authController.Register,
	)

	auth.POST("/login",
		middlewares.RateLimitPerRouteAndIP(1, 10*time.Second, 5),
		authController.Login,
	)

	auth.POST("/refresh-token",
		middlewares.RateLimitPerRouteAndIP(1, 5*time.Second, 10),
		authController.RefreshToken,
	)

	auth.POST("/logout",
		middlewares.RateLimitPerRouteAndIP(1, 5*time.Second, 5),
		authController.Logout,
	)

	auth.POST("/send-email/:email",
		middlewares.RateLimitPerRouteAndIP(1, 60*time.Second, 2),
		authController.SendVerifyEmail,
	)

	auth.GET("/verify-email/:token",
		middlewares.RateLimitPerRouteAndIP(1, 60*time.Second, 2),
		authController.VerifyEmail,
	)

	auth.POST("/reset-password/:email",
		middlewares.RateLimitPerRouteAndIP(1, 60*time.Second, 2),
		authController.SendResetPassword,
	)

	auth.POST("/change-password/:token",
		middlewares.RateLimitPerRouteAndIP(1, 20*time.Second, 3),
		authController.ResetPassword,
	)
}
