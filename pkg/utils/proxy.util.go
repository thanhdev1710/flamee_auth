package utils

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func ForwardTo(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
