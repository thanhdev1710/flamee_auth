package utils

import (
	"time"

	"github.com/gin-gonic/gin"
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

func SetCookiesToken(c *gin.Context, accessToken, refreshToken string, timeDefault, timeRemember time.Duration) {
	// Thiết lập cookie cho access token
	c.SetCookie("access_token", accessToken, int(timeDefault.Seconds()), "/", "localhost", true, true)

	// Thiết lập cookie cho refresh token
	c.SetCookie("refresh_token", refreshToken, int(timeRemember.Seconds()), "/", "localhost", true, true)
}
