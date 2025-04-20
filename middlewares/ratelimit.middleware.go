package middlewares

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"github.com/patrickmn/go-cache"
)

var buckets = cache.New(10*time.Minute, 1*time.Minute)

func RateLimitPerRouteAndIP(refillAmount int64, refillInterval time.Duration, capacity int64) gin.HandlerFunc {
	if refillAmount <= 0 {
		refillAmount = 1
	}
	if refillInterval <= 0 {
		refillInterval = time.Second
	}
	if capacity <= 0 {
		capacity = 10
	}

	// Tính khoảng thời gian hồi 1 token
	fillInterval := refillInterval / time.Duration(refillAmount)

	return func(c *gin.Context) {
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		ip := c.ClientIP()
		key := path + "_" + ip

		val, found := buckets.Get(key)
		if !found {
			bucket := ratelimit.NewBucket(fillInterval, capacity)
			buckets.Set(key, bucket, cache.DefaultExpiration)
			val = bucket
		}

		bucket := val.(*ratelimit.Bucket)

		if bucket.TakeAvailable(1) == 0 {
			c.Header("Retry-After", strconv.Itoa(int(fillInterval.Seconds())))
			c.JSON(429, gin.H{
				"status":  "error",
				"message": "Quá nhiều yêu cầu",
			})
			c.Abort()
			return
		}

		c.Header("X-RateLimit-Remaining", strconv.FormatInt(bucket.Available(), 10))
		c.Header("X-RateLimit-Limit", strconv.FormatInt(capacity, 10))
		c.Next()
	}
}
