package services

import (
	"crypto/tls"
	"fmt"
	"html"
	"time"

	"github.com/thanhdev1710/flamee_auth/global"
	gomail "gopkg.in/gomail.v2"
)

type EmailServices struct{}

func NewEmailServices() *EmailServices { return &EmailServices{} }

func (es *EmailServices) Send(toEmail, actionURL, emailType string) {
	from := global.Config.Email.Username     // mailbox thật trong Zoho
	password := global.Config.Email.Password // App Password nếu bật MFA
	host := global.Config.Email.Host         // "smtp.zoho.com"
	port := 465

	appName := "Flamee"
	subject := subjectFor(emailType, appName)

	// Nội dung ngắn gọn theo loại
	title, desc, cta := contentFor(emailType, appName)

	// Sanitize tối thiểu cho HTML
	action := html.EscapeString(actionURL)

	// --- HTML tối giản, table + inline style ---
	htmlBody := fmt.Sprintf(`<!doctype html>
<html><body style="margin:0;padding:0;background:#f5f5f7;">
  <!-- preheader ẩn -->
  <div style="display:none;opacity:0;visibility:hidden;mso-hide:all;line-height:1px;max-height:0;max-width:0;overflow:hidden;">%s</div>

  <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="background:#f5f5f7;padding:16px 0;">
    <tr><td align="center">
      <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="max-width:560px;background:#ffffff;border-radius:8px;padding:24px;">
        <tr><td align="left" style="font-family:Arial,Helvetica,sans-serif;color:#111;">
          <h1 style="font-size:18px;margin:0 0 8px;">%s</h1>
          <p style="font-size:14px;line-height:1.6;margin:0 0 16px;color:#444;">%s</p>
          <table role="presentation" cellspacing="0" cellpadding="0" style="margin:20px 0;">
            <tr>
              <td align="center" bgcolor="#4F46E5" style="border-radius:6px;">
                <a href="%s" target="_blank" style="display:inline-block;padding:12px 18px;color:#ffffff;text-decoration:none;font-weight:bold;font-family:Arial,Helvetica,sans-serif;">%s</a>
              </td>
            </tr>
          </table>
          <p style="font-size:12px;line-height:1.6;margin:16px 0 0;color:#666;">
            Nếu nút không hoạt động, hãy copy đường link bên dưới và dán vào trình duyệt:<br>
            <a href="%s" style="color:#4F46E5;text-decoration:underline;">%s</a>
          </p>
          <p style="font-size:12px;color:#999;margin:24px 0 0;">© %d %s</p>
        </td></tr>
      </table>
    </td></tr>
  </table>
</body></html>`,
		preheaderFor(emailType, appName),
		title, desc, action, cta, action, action, time.Now().Year(), appName,
	)

	// --- Text/plain fallback (rất quan trọng cho Gmail/Outlook) ---
	textBody := fmt.Sprintf("%s\n\n%s\n\nLink: %s\n", title, desc, actionURL)

	// Gửi bằng gomail
	m := gomail.NewMessage()
	m.SetHeader("From", from) // phải trùng mailbox AUTH
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)
	m.AddAlternative("text/plain", textBody)

	d := gomail.NewDialer(host, port, from, password)
	d.TLSConfig = &tls.Config{ServerName: host}

	go func() {
		if err := d.DialAndSend(m); err != nil {
			fmt.Printf("[email] ❌ send fail to %s: %v\n", toEmail, err)
		} else {
			fmt.Printf("[email] ✅ sent to %s | subject=%q\n", toEmail, subject)
		}
	}()
}

func subjectFor(emailType, app string) string {
	switch emailType {
	case "verification":
		return "Verify your email – " + app
	case "password_reset":
		return "Reset your password – " + app
	default:
		return app + " notification"
	}
}

func preheaderFor(emailType, app string) string {
	switch emailType {
	case "verification":
		return "Xác nhận email để kích hoạt tài khoản " + app + "."
	case "password_reset":
		return "Bạn vừa yêu cầu đặt lại mật khẩu."
	default:
		return "Thông báo từ " + app
	}
}

func contentFor(emailType, app string) (title, desc, cta string) {
	switch emailType {
	case "verification":
		return "Xác nhận email", "Nhấn nút bên dưới để xác nhận và kích hoạt tài khoản " + app + ".", "Xác nhận email"
	case "password_reset":
		return "Đặt lại mật khẩu", "Nhấn nút bên dưới để tiếp tục đặt lại mật khẩu cho tài khoản " + app + ".", "Đặt lại mật khẩu"
	default:
		return app, "Vui lòng mở liên kết bên dưới.", "Mở liên kết"
	}
}
