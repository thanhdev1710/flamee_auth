package profile

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/middlewares"
	"github.com/thanhdev1710/flamee_auth/pkg/utils"
)

type ProfileRouter struct{}

func (pr *ProfileRouter) InitProfileRouter(Router *gin.RouterGroup) {
	ProfileRouterPublic := Router.Group("/profiles")
	{
		ProfileRouterPublic.POST("/check-friend", utils.ForwardTo(global.Url.UrlUserService))
	}

	ProfileRouterPrivate := Router.Group("/profiles").
		Use(middlewares.AuthMiddleware()).
		Use(middlewares.VerifyAccount())
	{
		ProfileRouterPrivate.GET("", utils.ForwardTo(global.Url.UrlUserService))
		ProfileRouterPrivate.POST("", utils.ForwardTo(global.Url.UrlUserService))
		ProfileRouterPrivate.PUT("", utils.ForwardTo(global.Url.UrlUserService))

		ProfileRouterPrivate.GET("/online", utils.ForwardTo(global.Url.UrlUserService))
		ProfileRouterPrivate.GET("/search", utils.ForwardTo(global.Url.UrlUserService))
		ProfileRouterPrivate.GET("/search/:keyword", utils.ForwardTo(global.Url.UrlUserService))

		ProfileRouterPrivate.GET("/suggest-username/:base", utils.ForwardTo(global.Url.UrlUserService))
	}
}
