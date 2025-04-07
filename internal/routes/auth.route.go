package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/internal/controllers"
)

func AuthRoutes(r *gin.Engine) {
	auth := r.Group("auth")
	{
		auth.POST("/register", controllers.NewAuthControllers().Register)
		auth.POST("/login", controllers.NewAuthControllers().Login)
		auth.POST("/refresh-token", controllers.NewAuthControllers().RefreshToken)
		auth.POST("/logout", controllers.NewAuthControllers().Logout)
		auth.POST("/send-email/:email", controllers.NewAuthControllers().SendVerifyEmail)
		auth.GET("/verify-email/:token", controllers.NewAuthControllers().VerifyEmail)
		auth.POST("/reset-password/:email", controllers.NewAuthControllers().SendResetPassword)
		auth.POST("/change-password/:token", controllers.NewAuthControllers().ResetPassword)
	}
}
