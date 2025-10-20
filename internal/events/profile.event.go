package events

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/thanhdev1710/flamee_auth/internal/services"
)

type ProfileCreatedEvent struct {
	UserID string `json:"user_id"`
}

// HandleProfileCreated xử lý sự kiện profile.created
func HandleProfileCreated(m *nats.Msg) {
	var event ProfileCreatedEvent
	if err := json.Unmarshal(m.Data, &event); err != nil {
		log.Println("Lỗi khi giải mã dữ liệu:", err)
		return
	}

	if event.UserID == "" {
		log.Println("Lỗi: user_id không hợp lệ trong sự kiện profile.created")
		return
	}

	fmt.Printf("🔔 Sự kiện profile.created nhận được cho user_id: %s\n", event.UserID)
	services.NewUserServices().ConfirmProfile(event.UserID)
}
