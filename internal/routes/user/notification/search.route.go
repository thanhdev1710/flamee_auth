package notification

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/middlewares"
	"github.com/thanhdev1710/flamee_auth/pkg/utils"
)

type NotificationRouter struct{}

func (fr *NotificationRouter) InitNotificationRouter(Router *gin.RouterGroup) {
	NotificationRouterPrivate := Router.Group("/notifications").
		Use(middlewares.AuthMiddleware()).
		Use(middlewares.VerifyAccount())
	{
		NotificationRouterPrivate.POST("/", utils.ForwardTo(global.Url.UrlNotificationsService))
		NotificationRouterPrivate.GET("/", utils.ForwardTo(global.Url.UrlNotificationsService))
		NotificationRouterPrivate.PATCH("/:id/read", utils.ForwardTo(global.Url.UrlNotificationsService))
		NotificationRouterPrivate.PATCH("/read-all", utils.ForwardTo(global.Url.UrlNotificationsService))
		NotificationRouterPrivate.DELETE("/:id", utils.ForwardTo(global.Url.UrlNotificationsService))
	}
}
