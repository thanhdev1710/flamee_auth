package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/internal/repo"
	"github.com/thanhdev1710/flamee_auth/middlewares"
)

func ProtectRoutes(r *gin.Engine) {
	protect := r.Group("api/v1").Use(middlewares.AuthMiddleware())
	{
		protect.GET("/me", func(ctx *gin.Context) {
			// Lấy userId từ context
			userId := ctx.GetString("userId")
			if userId == "" {
				// Nếu không tìm thấy userId trong context, trả về lỗi
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "User ID not found in context",
				})
				return
			}

			// Tìm người dùng trong cơ sở dữ liệu theo userId
			user, err := repo.NewUserRepo().FindById(userId)
			if err != nil {
				// Nếu không tìm thấy người dùng hoặc có lỗi trong quá trình tìm kiếm, trả về lỗi
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			// Trả về thông tin người dùng
			ctx.JSON(http.StatusOK, user)

		})

		// 	protect.GET("/users/*action", func(ctx *gin.Context) {
		// 		proxy := utils.NewReverseProxy("http://localhost:8001") // Dịch vụ người dùng
		// 		proxy.ServeHTTP(ctx.Writer, ctx.Request)
		// })
	}
}
