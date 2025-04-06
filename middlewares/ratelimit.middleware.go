package middlewares

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

func RateLimitMiddleware(bucket *ratelimit.Bucket) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Kiểm tra và lấy số lượng yêu cầu còn lại
		remainingRequests := bucket.Available()

		// Nếu không còn yêu cầu, trả về lỗi 429
		if remainingRequests <= 0 {
			c.JSON(429, gin.H{
				"message":           "Too many requests",
				"remainingRequests": remainingRequests,
			})
			c.Abort()
			return
		}

		// Cập nhật số lượng yêu cầu còn lại
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(remainingRequests, 10))

		// Giảm số lượng yêu cầu còn lại sau mỗi request
		bucket.Take(1)

		// Tiến hành tiếp tục với request
		c.Next()
	}
}
