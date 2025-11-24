package profile

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/middlewares"
	"github.com/thanhdev1710/flamee_auth/pkg/utils"
)

type FollowRouter struct{}

func (fr *FollowRouter) InitFollowRouter(Router *gin.RouterGroup) {
	FollowRouterPublic := Router.Group("/follows")
	{
		FollowRouterPublic.POST("/check-friend", utils.ForwardTo(global.Url.UrlUserService))
	}

	FollowRouterPrivate := Router.Group("/follows").
		Use(middlewares.AuthMiddleware()).
		Use(middlewares.VerifyAccount())
	{
		FollowRouterPrivate.GET("/friend_suggestions", utils.ForwardTo(global.Url.UrlUserService))
		FollowRouterPrivate.GET("/friend_suggestions/:username", utils.ForwardTo(global.Url.UrlUserService))
		FollowRouterPrivate.POST("", utils.ForwardTo(global.Url.UrlUserService))
	}
}
