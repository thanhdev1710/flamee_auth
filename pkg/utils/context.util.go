package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/global"
)

// GetUserId lấy thông tin userId từ context
func GetUserId(c *gin.Context) string {
	userId, _ := c.Get("userId")
	if userId != nil {
		return userId.(string)
	}
	return ""
}

// GetRole lấy thông tin role từ context
func GetRole(c *gin.Context) string {
	role, _ := c.Get("role")
	if role != nil {
		return role.(string)
	}
	return ""
}

func SetCookiesToken(w http.ResponseWriter, accessToken, refreshToken string, timeDefault, timeRemember time.Duration) {
	// ⚠️ Thay bằng domain backend thật sự khi deploy

	http.SetCookie(w, &http.Cookie{
		Name:     HexString("flamee_access_token"),
		Value:    accessToken,
		Path:     "/",
		Domain:   global.Config.Domain,
		MaxAge:   int(timeDefault.Seconds()),
		HttpOnly: true,
		Secure:   true,                  // ✅ Bắt buộc khi dùng HTTPS
		SameSite: http.SameSiteNoneMode, // ✅ Cho phép cookie cross-domain
	})

	http.SetCookie(w, &http.Cookie{
		Name:     HexString("flamee_refresh_token"),
		Value:    refreshToken,
		Path:     "/",
		Domain:   global.Config.Domain,
		MaxAge:   int(timeRemember.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
}
