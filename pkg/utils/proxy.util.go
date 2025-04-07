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
		userId := c.GetString("userId")
		role := c.GetString("role")

		// Kiểm tra nếu không tìm thấy userId hoặc role trong context
		if userId == "" || role == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID or Role not found in context"})
			return
		}

		remote, err := url.Parse(target)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid target URL"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)

		// Cập nhật scheme và host, giữ nguyên Path và Query gốc
		c.Request.URL.Scheme = remote.Scheme
		c.Request.URL.Host = remote.Host

		// Optionally: chỉnh lại Host header cho đúng server đích (nếu phía sau kiểm tra host)
		c.Request.Host = remote.Host

		// Thêm userId và role vào header yêu cầu
		c.Request.Header.Add("userId", userId)
		c.Request.Header.Add("role", role)

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
