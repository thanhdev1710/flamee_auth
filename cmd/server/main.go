package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/thanhdev1710/flamee_auth/internal/initialize"
)

func main() {
	// Khởi tạo ứng dụng
	initialize.Run()

	// Đảm bảo đóng kết nối khi ứng dụng nhận tín hiệu dừng (Ctrl+C, SIGTERM)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Đợi tín hiệu dừng
	<-signalChan
	log.Println("Đang dừng ứng dụng...")

	// Đảm bảo đóng tài nguyên khi ứng dụng thoát
	// Đóng kết nối NATS và các tài nguyên khác
	initialize.CloseNats()
	initialize.ClosePostgreSql()

	// Nếu có kết nối Redis hoặc PostgreSQL, bạn cũng có thể đóng chúng ở đây
	// initialize.CloseRedis()

	log.Println("Ứng dụng đã dừng thành công")
}
