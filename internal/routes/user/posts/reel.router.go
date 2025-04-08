package posts

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/middlewares"
	"github.com/thanhdev1710/flamee_auth/pkg/utils"
)

type ReelRouter struct{}

func (rr *ReelRouter) InitReelRouter(Router *gin.RouterGroup) {
	// Public router
	ReelRouterPublic := Router.Group("/reel")
	{
		ReelRouterPublic.GET("/:id", utils.ForwardTo(global.Url.UrlPostService))
	}
	// Private router
	ReelRouterPrivate := Router.Group("/reel").
		Use(middlewares.AuthMiddleware()).
		Use(middlewares.VerifyAccount())
	{
		ReelRouterPrivate.POST("/", utils.ForwardTo(global.Url.UrlPostService))
		ReelRouterPrivate.DELETE("/:id", utils.ForwardTo(global.Url.UrlPostService))
	}
}
