package services

import (
	"fmt"
	"net/smtp"

	"github.com/thanhdev1710/flamee_auth/global"
)

type EmailServices struct{}

func NewEmailServices() *EmailServices {
	return &EmailServices{}
}

func (es *EmailServices) SendVerificationEmail(toEmail string, verificationURL string) {
	// Cấu hình người gửi
	from := global.Config.Email.Username
	password := global.Config.Email.Password

	// Cấu hình SMTP server
	smtpHost := global.Config.Email.Host
	smtpPort := global.Config.Email.Port

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
			// Nếu gửi không thành công, có thể ghi log lỗi ở đây
			fmt.Printf("Failed to send verification email to %s: %v\n", toEmail, err)
		} else {
			// Nếu gửi thành công, ghi log thành công (nếu cần)
			fmt.Printf("Verification email sent to %s\n", toEmail)
		}
	}()
}
