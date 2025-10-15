package initialize

import (
	"fmt"

	"github.com/thanhdev1710/flamee_auth/global"
)

func Run() {
	LoadConfig()
	InitLogger()
	InitPostgreSql()
	InitNats()
	r := InitRouter()

	fmt.Print(global.Url)

	port := global.Config.Port
	if port == "" {
		port = "8081"
	}

	r.Run(":" + port)
}
