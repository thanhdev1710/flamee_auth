package user

import (
	"github.com/thanhdev1710/flamee_auth/internal/routes/user/posts"
)

type RouterGroup struct {
	Post posts.RouterGroup
}
