package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/internal/repo"
	"github.com/thanhdev1710/flamee_auth/middlewares"
	"github.com/thanhdev1710/flamee_auth/pkg/utils"
)

func ProtectRoutes(r *gin.Engine) {
	protect := r.Group("api/v1").Use(middlewares.AuthMiddleware()).Use(middlewares.VerifyEmail())
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

		protect.POST("/interactions/like/:postId", utils.ForwardTo("http://localhost:3000"))
		protect.DELETE("/interactions/like/:postId", utils.ForwardTo("http://localhost:3000"))
		protect.POST("/interactions/comment/:postId", utils.ForwardTo("http://localhost:3000"))
		protect.DELETE("/interactions/comment/:postId", utils.ForwardTo("http://localhost:3000"))
		protect.POST("/interactions/share/:postId", utils.ForwardTo("http://localhost:3000"))
		protect.GET("/interactions/:postId", utils.ForwardTo("http://localhost:3000"))

		protect.POST("/posts", utils.ForwardTo("http://localhost:3000"))
		protect.PUT("/posts/:id", utils.ForwardTo("http://localhost:3000"))
		protect.DELETE("/posts/:id", utils.ForwardTo("http://localhost:3000"))
		protect.GET("/posts/:id", utils.ForwardTo("http://localhost:3000"))
	}
}
