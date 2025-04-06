package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/middlewares"
)

func ProtectRoutes(r *gin.Engine) {
	protect := r.Group("api/v1").Use(middlewares.AuthMiddleware())
	{
		protect.GET("/me", func(ctx *gin.Context) {
			userId := ctx.GetString("userId")
			role := ctx.GetString("role")
			ctx.JSON(200, gin.H{
				"userId": userId,
				"role":   role,
			})
		})

		// 	protect.GET("/users/*action", func(ctx *gin.Context) {
		// 		proxy := utils.NewReverseProxy("http://localhost:8001") // Dịch vụ người dùng
		// 		proxy.ServeHTTP(ctx.Writer, ctx.Request)
		// })
	}
}
