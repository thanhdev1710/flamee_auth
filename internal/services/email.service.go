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

func (es *EmailServices) Send(toEmail string, verificationURL string, emailType string) {
	// Cấu hình người gửi
	from := global.Config.Email.Username
	password := global.Config.Email.Password

	// Cấu hình SMTP server
	smtpHost := global.Config.Email.Host
	smtpPort := global.Config.Email.Port

	// Tạo nội dung email HTML
	var htmlBody string
	var subject string

	switch emailType {
	case "verification":
		// Nội dung email xác thực email
		htmlBody = "<html><body>" +
			"<h3>Verify your email</h3>" +
			"<p>Click the button below to verify your email address:</p>" +
			"<a href=\"" + verificationURL + "\" style=\"background-color: #4CAF50; color: white; padding: 14px 20px; text-align: center; text-decoration: none; display: inline-block; border-radius: 4px;\">Verify Email</a>" +
			"</body></html>"
		subject = "Subject: Email Verification\r\n"

	case "password_reset":
		// Nội dung email reset password
		htmlBody = "<html><body>" +
			"<h3>Password Reset Request</h3>" +
			"<p>Click the button below to reset your password:</p>" +
			"<a href=\"" + verificationURL + "\" style=\"background-color: #FF6347; color: white; padding: 14px 20px; text-align: center; text-decoration: none; display: inline-block; border-radius: 4px;\">Reset Password</a>" +
			"</body></html>"
		subject = "Subject: Password Reset Request\r\n"

	default:
		// Mặc định nếu không có loại email nào khớp
		htmlBody = "<html><body>" +
			"<h3>Unknown Email Type</h3>" +
			"<p>We are sorry, but we could not process your request.</p>" +
			"</body></html>"
		subject = "Subject: Unknown Request\r\n"
	}

	// MIME header và nội dung body
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
			fmt.Printf("Failed to send to %s: %v\n", toEmail, err)
		} else {
			// Nếu gửi thành công, ghi log thành công (nếu cần)
			fmt.Printf("Email sent to %s\n", toEmail)
		}
	}()
}
