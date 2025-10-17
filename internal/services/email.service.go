package services

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/thanhdev1710/flamee_auth/global"
	gomail "gopkg.in/gomail.v2"
)

type EmailServices struct{}

func NewEmailServices() *EmailServices { return &EmailServices{} }

type emailData struct {
	AppName      string
	BrandColor   string
	LogoURL      string
	ActionURL    string
	SupportEmail string
	Recipient    string
	Year         int
	EmailType    string // "verification" | "password_reset" | ...
}

// Send: gửi email đẹp (HTML + text) qua Zoho SMTPS 465
func (es *EmailServices) Send(toEmail string, actionURL string, emailType string) {
	from := global.Config.Email.Username     // ví dụ: no-reply@yourdomain.com (phải trùng mailbox AUTH)
	password := global.Config.Email.Password // App password nếu bật MFA
	smtpHost := global.Config.Email.Host     // "smtp.zoho.com"
	smtpPort := 465

	// === Branding / cấu hình hiển thị ===
	data := emailData{
		AppName:      "Flamee",
		BrandColor:   "#4F46E5",                                                 // Indigo-600
		LogoURL:      "https://dummyimage.com/140x40/4F46E5/ffffff&text=Flamee", // thay bằng logo thực
		ActionURL:    actionURL,
		SupportEmail: "support@yourdomain.com",
		Recipient:    toEmail,
		Year:         time.Now().Year(),
		EmailType:    emailType,
	}

	subject := subjectFor(emailType, data.AppName)
	htmlBody, textBody, err := renderEmail(emailType, data)
	if err != nil {
		fmt.Printf("[email] ❌ template render error: %v\n", err)
		return
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)

	// Preheader ẩn (giúp tăng CTR, hiển thị trong inbox preview)
	preheader := preheaderFor(emailType, data.AppName)
	// chèn preheader (ẩn) vào đầu body HTML
	htmlBody = fmt.Sprintf(`<div style="display:none;opacity:0;visibility:hidden;mso-hide:all;line-height:1px;max-height:0;max-width:0;overflow:hidden;">%s</div>%s`, template.HTMLEscapeString(preheader), htmlBody)

	m.SetBody("text/html", htmlBody)
	m.AddAlternative("text/plain", textBody)

	d := gomail.NewDialer(smtpHost, smtpPort, from, password)
	d.TLSConfig = &tls.Config{ServerName: smtpHost}

	// Gửi đồng bộ; nếu muốn async: bọc go func(){...}()
	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("[email] ❌ Failed to send to %s: %v\n", toEmail, err)
	} else {
		fmt.Printf("[email] ✅ Email sent successfully to %s | subject=%q\n", toEmail, subject)
	}
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
		return "Confirm your email to activate your " + app + " account."
	case "password_reset":
		return "We received a request to reset your password."
	default:
		return "A message from " + app
	}
}

// ====== TEMPLATE ======

var tpl = template.Must(template.New("email").Parse(strings.TrimSpace(`
{{- /* HTML email template */ -}}
<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="x-apple-disable-message-reformatting">
  <meta name="color-scheme" content="light dark">
  <meta name="supported-color-schemes" content="light dark">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{.AppName}}</title>
  <style>
    /* Reset cơ bản cho email client */
    body, table, td, a { font-family: Inter, system-ui, -apple-system, Segoe UI, Roboto, sans-serif; }
    img { border:0; outline:none; text-decoration:none; }
    table { border-collapse: collapse !important; }
    a { text-decoration: none; }
    /* Container */
    .wrapper { width:100%; background:#f6f7fb; margin:0; padding:24px; }
    .card { width:100%; max-width:560px; margin:0 auto; background:#ffffff; border-radius:16px; padding:32px; box-shadow: 0 2px 8px rgba(0,0,0,0.05); }
    .logo { text-align:center; margin-bottom:16px; }
    .title { font-size:20px; font-weight:700; color:#0f172a; margin:0 0 12px; text-align:center; }
    .desc { font-size:14px; color:#334155; line-height:1.6; margin:0 0 20px; text-align:center; }
    .cta-wrap { text-align:center; margin:24px 0 12px; }
    .btn { display:inline-block; padding:12px 20px; border-radius:10px; background:{{.BrandColor}}; color:#ffffff !important; font-weight:600; }
    .muted { font-size:12px; color:#64748b; text-align:center; margin:20px 0 0; }
    .footer { font-size:12px; color:#94a3b8; text-align:center; margin-top:24px; }
    .link { color:{{.BrandColor}} !important; text-decoration:underline; }
    /* Dark mode */
    @media (prefers-color-scheme: dark) {
      body, .wrapper { background:#0b1020 !important; }
      .card { background:#121a2b !important; box-shadow:none; }
      .title { color:#e2e8f0 !important; }
      .desc { color:#cbd5e1 !important; }
      .muted { color:#94a3b8 !important; }
      .footer { color:#64748b !important; }
    }
    /* Mobile */
    @media screen and (max-width: 600px) {
      .card { padding:24px; border-radius:14px; }
      .title { font-size:18px; }
    }
  </style>
</head>
<body>
  <div class="wrapper">
    <table role="presentation" class="card" cellspacing="0" cellpadding="0">
      <tr>
        <td>
          <div class="logo">
            {{if .LogoURL}}<img src="{{.LogoURL}}" width="140" height="40" alt="{{.AppName}} logo" loading="lazy">{{else}}<h2 style="margin:0">{{.AppName}}</h2>{{end}}
          </div>

          {{if eq .EmailType "verification"}}
            <h1 class="title">Verify your email</h1>
            <p class="desc">Thanks for signing up for <strong>{{.AppName}}</strong>! Please confirm your email to activate your account.</p>
          {{else if eq .EmailType "password_reset"}}
            <h1 class="title">Reset your password</h1>
            <p class="desc">You requested to reset your password for <strong>{{.AppName}}</strong>. Click the button below to continue.</p>
          {{else}}
            <h1 class="title">{{.AppName}}</h1>
            <p class="desc">We have an update for your account.</p>
          {{end}}

          <div class="cta-wrap">
            <a class="btn" href="{{.ActionURL}}" target="_blank" rel="noopener">Open Link</a>
          </div>

          <p class="muted">
            If the button doesn’t work, copy and paste this link into your browser:<br>
            <a class="link" href="{{.ActionURL}}" target="_blank" rel="noopener">{{.ActionURL}}</a>
          </p>

          <p class="footer">
            Need help? Contact us at <a class="link" href="mailto:{{.SupportEmail}}">{{.SupportEmail}}</a><br>
            © {{.Year}} {{.AppName}}. All rights reserved.
          </p>
        </td>
      </tr>
    </table>
  </div>
</body>
</html>
`)))

func renderEmail(emailType string, d emailData) (htmlOut string, textOut string, err error) {
	var sb strings.Builder
	if err = tpl.Execute(&sb, d); err != nil {
		return "", "", err
	}
	htmlOut = sb.String()

	// Bản text (fallback)
	switch emailType {
	case "verification":
		textOut = fmt.Sprintf(
			"Verify your email for %s\n\nOpen this link to confirm:\n%s\n\nIf you didn’t request this, you can ignore this email.\n",
			d.AppName, d.ActionURL,
		)
	case "password_reset":
		textOut = fmt.Sprintf(
			"Reset your password for %s\n\nOpen this link to continue:\n%s\n\nIf you didn’t request this, you can ignore this email.\n",
			d.AppName, d.ActionURL,
		)
	default:
		textOut = fmt.Sprintf(
			"%s notification\n\nOpen this link:\n%s\n", d.AppName, d.ActionURL,
		)
	}
	return htmlOut, textOut, nil
}
