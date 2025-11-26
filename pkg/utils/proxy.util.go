package utils

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func ForwardTo(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy JWT từ context
		jwtToken := c.GetString("jwt")

		// Parse URL đích
		remote, err := url.Parse(target)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "URL đích không hợp lệ",
			})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)

		// Cập nhật scheme và host
		c.Request.URL.Scheme = remote.Scheme
		c.Request.URL.Host = remote.Host
		c.Request.Host = remote.Host // nếu cần kiểm tra Host header phía sau
		// Gắn lại Authorization header

		if jwtToken != "" {
			c.Request.Header.Set("Authorization", jwtToken)
		}
		// Forward request
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
