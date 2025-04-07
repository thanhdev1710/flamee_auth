package middlewares

import "github.com/gin-gonic/gin"

func Helmet(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Next()
	})
}
