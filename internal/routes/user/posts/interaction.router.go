package posts

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/middlewares"
	"github.com/thanhdev1710/flamee_auth/pkg/utils"
)

type InteractionRouter struct{}

func (ir *InteractionRouter) InitInteractionRouter(Router *gin.RouterGroup) {
	// Public router
	InteractionRouterPublic := Router.Group("/interactions")
	{
		InteractionRouterPublic.GET("/:postId", utils.ForwardTo("http://localhost:3000"))
	}
	// Private router
	InteractionRouterPrivate := Router.Group("/interactions").
		Use(middlewares.AuthMiddleware()).
		Use(middlewares.VerifyEmail())
	{
		InteractionRouterPrivate.POST("/like/:postId", utils.ForwardTo("http://localhost:3000"))
		InteractionRouterPrivate.DELETE("/like/:postId", utils.ForwardTo("http://localhost:3000"))
		InteractionRouterPrivate.POST("/comment/:postId", utils.ForwardTo("http://localhost:3000"))
		InteractionRouterPrivate.DELETE("/comment/:postId", utils.ForwardTo("http://localhost:3000"))
		InteractionRouterPrivate.POST("/share/:postId", utils.ForwardTo("http://localhost:3000"))
	}
}
