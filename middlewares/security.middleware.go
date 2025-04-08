package middlewares

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/global"
)

func CheckAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy API Key từ header
		apiKey := c.GetHeader("X-API-Key")

		// Kiểm tra API Key hợp lệ (sử dụng biến môi trường hoặc key hardcoded)
		if apiKey != global.Config.ApiKey { // Hoặc một giá trị cố định
			c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden"})
			c.Abort()
			return
		}

		// Nếu API Key hợp lệ, tiếp tục xử lý
		c.Next()
	}
}

// Middleware chặn các User-Agent đáng ngờ
func BlockSuspiciousUserAgents() gin.HandlerFunc {
	// Một vài regex pattern đơn giản (tối ưu theo danh sách dài bạn đưa)
	patterns := []string{
		`(?i)curl`, `(?i)wget`, `(?i)httpie`, `(?i)python`, `(?i)go-http-client`,
		`(?i)java`, `(?i)libwww-perl`, `(?i)lynx`, `(?i)scrapy`, `(?i)node`,
		`(?i)phantomjs`, `(?i)headlesschrome`, `(?i)sqlmap`, `(?i)nikto`, `(?i)nmap`,
		`(?i)zap`, `(?i)fuzz`, `(?i)bot`, `(?i)spider`,
		`(?i)crawler`, `(?i)masscan`, `(?i)scan`, `(?i)grab`, `(?i)fetch`,
	}

	// , `(?i)postmanruntime`

	var regexList []*regexp.Regexp
	for _, p := range patterns {
		re := regexp.MustCompile(p)
		regexList = append(regexList, re)
	}

	return func(c *gin.Context) {
		ua := c.GetHeader("User-Agent")
		if ua == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "Access denied: missing User-Agent",
			})
			return
		}

		for _, re := range regexList {
			if re.MatchString(ua) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"message": "Access denied: suspicious User-Agent",
				})
				return
			}
		}
		c.Next()
	}
}
