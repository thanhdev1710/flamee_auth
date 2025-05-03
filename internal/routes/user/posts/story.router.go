package posts

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/middlewares"
	"github.com/thanhdev1710/flamee_auth/pkg/utils"
)

type StoryRouter struct{}

func (sr *StoryRouter) InitStoryRouter(Router *gin.RouterGroup) {
	// Public router
	StoryRouterPublic := Router.Group("/story")
	{
		StoryRouterPublic.GET("/:id", utils.ForwardTo(global.Url.UrlPostService))
	}
	// Private router
	StoryRouterPrivate := Router.Group("/story").
		Use(middlewares.AuthMiddleware()).
		Use(middlewares.VerifyAccount())
	{
		StoryRouterPrivate.POST("", utils.ForwardTo(global.Url.UrlPostService))
		StoryRouterPrivate.DELETE("/:id", utils.ForwardTo(global.Url.UrlPostService))
	}
}
