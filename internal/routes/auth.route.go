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
		auth.POST("/refresh_token", controllers.NewAuthControllers().RefreshToken)
	}
}
