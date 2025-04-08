package routes

import (
	"github.com/thanhdev1710/flamee_auth/internal/routes/admin"
	"github.com/thanhdev1710/flamee_auth/internal/routes/auth"
	"github.com/thanhdev1710/flamee_auth/internal/routes/user"
)

type RouterGroup struct {
	User  user.RouterGroup
	Admin admin.RouterGroup
	Auth  auth.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
