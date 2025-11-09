package search

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/middlewares"
	"github.com/thanhdev1710/flamee_auth/pkg/utils"
)

type SearchRouter struct{}

func (fr *SearchRouter) InitSearchRouter(Router *gin.RouterGroup) {
	SearchRouterPrivate := Router.Group("/search").
		Use(middlewares.AuthMiddleware()).
		Use(middlewares.VerifyAccount())
	{
		SearchRouterPrivate.GET("/hot", utils.ForwardTo(global.Url.UrlSearchService))
		SearchRouterPrivate.GET("/users", utils.ForwardTo(global.Url.UrlSearchService))
		SearchRouterPrivate.GET("/posts", utils.ForwardTo(global.Url.UrlSearchService))
		SearchRouterPrivate.GET("/posts/hashtag", utils.ForwardTo(global.Url.UrlSearchService))
		// SearchRouterPrivate.GET("/admin/search/posts", utils.ForwardTo(global.Url.UrlSearchService))
		SearchRouterPrivate.GET("/posts/:id", utils.ForwardTo(global.Url.UrlSearchService))
	}
}
