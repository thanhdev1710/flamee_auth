package admin

import "github.com/thanhdev1710/flamee_auth/internal/routes/admin/posts"

type RouterGroup struct {
	Post posts.RouterGroup
}
