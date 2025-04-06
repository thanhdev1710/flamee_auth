package initialize

import (
	"github.com/thanhdev1710/flamee_auth/global"
)

func Run() {
	LoadConfig()
	InitPostgreSql()
	r := InitRouter()

	r.Run(":" + global.Config.Port)
}
