package posts

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/middlewares"
	"github.com/thanhdev1710/flamee_auth/pkg/utils"
)

type PostRouter struct{}

func (pr *PostRouter) InitPostRouter(Router *gin.RouterGroup) {
	// Public router
	PostRouterPublic := Router.Group("/admin/posts")
	{
		PostRouterPublic.GET("/")
		PostRouterPublic.GET("/:id", utils.ForwardTo("http://localhost:3000"))
	}
	// Private router
	PostRouterPrivate := Router.Group("/admin/posts").
		Use(middlewares.AuthMiddleware()).
		Use(middlewares.VerifyEmail())
	{
		PostRouterPrivate.POST("/", utils.ForwardTo("http://localhost:3000"))
		PostRouterPrivate.PUT("/:id", utils.ForwardTo("http://localhost:3000"))
		PostRouterPrivate.DELETE("/:id", utils.ForwardTo("http://localhost:3000"))
	}
}
