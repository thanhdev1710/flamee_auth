package services

import (
	"fmt"
	"log"
	"net/smtp"
)

type EmailServices struct{}

func NewEmailServices() *EmailServices {
	return &EmailServices{}
}

// TODO: MAI CODE TIẾP PHẦN GỬI EMAIL
func SendVerificationEmail(toEmail string, verificationURL string) error {
	// Cấu hình người gửi
	from := "youremail@example.com"
	password := "yourpassword"

	// Cấu hình SMTP server
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Tạo nội dung email HTML với button xác thực
	htmlBody := fmt.Sprintf(`
		<html>
		<body>
			<h3>Verify your email</h3>
			<p>Click the button below to verify your email address:</p>
			<a href="%s" style="background-color: #4CAF50; color: white; padding: 14px 20px; text-align: center; text-decoration: none; display: inline-block; border-radius: 4px;">Verify Email</a>
		</body>
		</html>
	`, verificationURL)

	// Cấu trúc email
	subject := "Subject: Email Verification\r\n"
	body := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n" + htmlBody
	msg := []byte(subject + body)

	// Kết nối với SMTP server
	auth := smtp.PlainAuth("", from, password, smtpHost)
	to := []string{toEmail}

	// Sử dụng go routine để gửi email bất đồng bộ
	go func() {
		// Gửi email
		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
		if err != nil {
			log.Println("Failed to send email:", err)
			return
		}

		// In ra khi email đã được gửi thành công
		fmt.Println("Verification email sent successfully!")
	}()

	// Trả về nil vì chúng ta không cần chờ đợi việc gửi email
	return nil
}
