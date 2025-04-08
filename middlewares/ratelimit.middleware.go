package middlewares

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"github.com/patrickmn/go-cache"
)

var buckets = cache.New(10*time.Minute, 1*time.Minute)

func RateLimitPerRouteAndIP(rate time.Duration, capacity int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		ip := c.ClientIP()
		key := path + "_" + ip

		val, found := buckets.Get(key)
		if !found {
			val = ratelimit.NewBucket(rate, capacity)
			buckets.Set(key, val, cache.DefaultExpiration)
		}

		bucket := val.(*ratelimit.Bucket)

		if bucket.TakeAvailable(1) == 0 {
			c.Header("Retry-After", strconv.Itoa(int(rate.Seconds())))
			c.JSON(429, gin.H{
				"message": "Too many requests",
			})
			c.Abort()
			return
		}

		c.Header("X-RateLimit-Remaining", strconv.FormatInt(bucket.Available(), 10))
		c.Header("X-RateLimit-Limit", strconv.FormatInt(capacity, 10))
		c.Next()
	}
}
