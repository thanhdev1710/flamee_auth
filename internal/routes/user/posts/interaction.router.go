package posts

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/middlewares"
	"github.com/thanhdev1710/flamee_auth/pkg/utils"
)

type InteractionRouter struct{}

func (ir *InteractionRouter) InitInteractionRouter(Router *gin.RouterGroup) {
	// Public router
	InteractionRouterPublic := Router.Group("/interactions")
	{
		InteractionRouterPublic.GET("/:postId", utils.ForwardTo(global.Url.UrlPostService))
	}
	// Private router
	InteractionRouterPrivate := Router.Group("/interactions").
		Use(middlewares.AuthMiddleware()).
		Use(middlewares.VerifyAccount())
	{
		InteractionRouterPrivate.POST("/like/:postId", utils.ForwardTo(global.Url.UrlPostService))
		InteractionRouterPrivate.DELETE("/like/:postId", utils.ForwardTo(global.Url.UrlPostService))
		InteractionRouterPrivate.POST("/comment/:postId", utils.ForwardTo(global.Url.UrlPostService))
		InteractionRouterPrivate.DELETE("/comment/:postId", utils.ForwardTo(global.Url.UrlPostService))
		InteractionRouterPrivate.POST("/share/:postId", utils.ForwardTo(global.Url.UrlPostService))
	}
}
