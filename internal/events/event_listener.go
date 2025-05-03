package events

import (
	"log"

	"github.com/thanhdev1710/flamee_auth/global"
)

// RegisterEventListeners đăng ký các sự kiện cần lắng nghe
func RegisterEventListeners() {
	// Lắng nghe sự kiện profile.created
	_, err := global.NatsConnection.Subscribe("profile.created", HandleProfileCreated)
	if err != nil {
		log.Fatalf("❌ Lỗi khi đăng ký sự kiện profile.created: %v", err)
	}

}
