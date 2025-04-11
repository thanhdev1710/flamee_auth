package profile

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/pkg/utils"
)

type FollowRouter struct{}

func (fr *FollowRouter) InitFollowRouter(Router *gin.RouterGroup) {
	FollowRouterPublic := Router.Group("/follows")
	{
		FollowRouterPublic.POST("/check-friend", utils.ForwardTo(global.Url.UrlUserService))
	}
}
