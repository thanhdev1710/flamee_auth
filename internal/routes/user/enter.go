package user

import (
	"github.com/thanhdev1710/flamee_auth/internal/routes/user/posts"
	"github.com/thanhdev1710/flamee_auth/internal/routes/user/profile"
)

type RouterGroup struct {
	Post    posts.RouterGroup
	Profile profile.RouterGroup
}
