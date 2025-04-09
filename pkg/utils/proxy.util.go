package utils

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func ForwardTo(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy userId và role từ context
		userId := GetUserId(c)
		role := c.GetString("role")

		// Kiểm tra nếu không tìm thấy userId hoặc role trong context
		if userId == "" || role == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Không tìm thấy User ID hoặc quyền trong context"})
			return
		}

		remote, err := url.Parse(target)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "URL đích không hợp lệ"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)

		// Cập nhật scheme và host, giữ nguyên Path và Query gốc
		c.Request.URL.Scheme = remote.Scheme
		c.Request.URL.Host = remote.Host

		// Tùy chọn: chỉnh lại Host header cho đúng server đích
		c.Request.Host = remote.Host

		// Thêm userId và role vào header yêu cầu
		c.Request.Header.Add("userId", userId)
		c.Request.Header.Add("role", role)

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
