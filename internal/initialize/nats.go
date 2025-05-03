package initialize

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/internal/events"
)

// InitNats khởi tạo kết nối NATS với các cấu hình chuẩn
func InitNats() {
	natsConfig := global.Config.Nats

	// Cấu hình NATS URL (có thể là IP hoặc hostname)
	natsURL := fmt.Sprintf("nats://%s:%s", natsConfig.Host, natsConfig.Port)

	// Cấu hình các tham số kết nối lại và tối đa các lần kết nối lại
	opts := []nats.Option{
		nats.MaxReconnects(natsConfig.MaxReconnects),
		nats.ReconnectWait(time.Duration(natsConfig.ReconnectWait) * time.Second),
		nats.PingInterval(30 * time.Second),
	}

	// Kết nối tới NATS
	nc, err := nats.Connect(natsURL, opts...)
	if err != nil {
		log.Fatalf("❌ Không thể kết nối đến NATS: %v", err)
	}

	// Lưu kết nối NATS vào biến toàn cục
	global.NatsConnection = nc
	// Kiểm tra trạng thái kết nối
	if global.NatsConnection != nil {
		log.Println("✅ Kết nối NATS thành công")
	}

	// Đăng ký các sự kiện
	events.RegisterEventListeners()

}

// CloseNats đóng kết nối NATS khi không còn sử dụng
func CloseNats() {
	if global.NatsConnection != nil {
		// Đảm bảo đóng kết nối NATS
		global.NatsConnection.Close()
		log.Println("✅ Đóng kết nối NATS")
	}
}
