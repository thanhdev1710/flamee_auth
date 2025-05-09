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

func SetCookiesToken(c *gin.Context, accessToken, refreshToken *string, timeDefault, timeRemember time.Duration) {
	c.SetSameSite(http.SameSiteNoneMode)

	// Cookie Access Token
	if accessToken != nil {
		c.SetCookie(
			HexString(global.Token.AccessToken), // Tên cookie
			*accessToken,                        // Giá trị
			int(timeDefault.Seconds()),          // Thời gian sống (giây)
			"",                                  // Path
			global.Config.Domain,                // Domain (ví dụ: flamee-auth.onrender.com)
			true,                                // Secure (chỉ gửi qua HTTPS)
			true,                                // HttpOnly
		)
	}

	// Cookie Refresh Token
	if refreshToken != nil {
		c.SetCookie(
			HexString(global.Token.RefreshToken),
			*refreshToken,
			int(timeRemember.Seconds()),
			"",
			global.Config.Domain,
			true,
			true,
		)
	}
}
