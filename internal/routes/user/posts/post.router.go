package posts

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/middlewares"
	"github.com/thanhdev1710/flamee_auth/pkg/utils"
)

type PostRouter struct{}

func (pr *PostRouter) InitPostRouter(Router *gin.RouterGroup) {
	// Public router
	PostRouterPublic := Router.Group("/posts")
	{
		PostRouterPublic.GET("/")
		PostRouterPublic.GET("/:id", utils.ForwardTo(global.Url.UrlPostService))
	}
	// Private router
	PostRouterPrivate := Router.Group("/posts").
		Use(middlewares.AuthMiddleware()).
		Use(middlewares.VerifyAccount())
	{
		PostRouterPrivate.POST("/", utils.ForwardTo(global.Url.UrlPostService))
		PostRouterPrivate.PUT("/:id", utils.ForwardTo(global.Url.UrlPostService))
		PostRouterPrivate.DELETE("/:id", utils.ForwardTo(global.Url.UrlPostService))
	}
}
