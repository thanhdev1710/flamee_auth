package user

import (
	"github.com/thanhdev1710/flamee_auth/internal/routes/user/notification"
	"github.com/thanhdev1710/flamee_auth/internal/routes/user/posts"
	"github.com/thanhdev1710/flamee_auth/internal/routes/user/profile"
	"github.com/thanhdev1710/flamee_auth/internal/routes/user/search"
)

type RouterGroup struct {
	Post         posts.RouterGroup
	Profile      profile.RouterGroup
	Search       search.RouterGroup
	Notification notification.RouterGroup
}
