package log

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/global"
)

type LogEntry struct {
	Level     string  `json:"level"`
	Time      string  `json:"time"`
	Caller    string  `json:"caller"`
	Message   string  `json:"msg"`
	Status    int     `json:"status"`
	Method    string  `json:"method"`
	Path      string  `json:"path"`
	ClientIP  string  `json:"clientIP"`
	UserAgent string  `json:"userAgent"`
	UserId    string  `json:"userId"`
	Latency   float64 `json:"latency"`
}

type LogRouter struct {
}

func (lr *LogRouter) InitLogRouter(r *gin.RouterGroup) {
	r.GET("/logs", getLogs)
}

func getLogs(c *gin.Context) {
	fmt.Println(global.Config.Logger.Filename)
	file, err := os.Open(global.Config.Logger.Filename)
	if err != nil {
		c.JSON(500, gin.H{"error": "cannot read log file"})
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	logs := make([]LogEntry, 0)

	// Query filter
	level := c.Query("level")
	method := c.Query("method")
	path := c.Query("path")
	status := c.Query("status")

	for scanner.Scan() {
		var entry LogEntry
		line := scanner.Bytes()

		if err := json.Unmarshal(line, &entry); err != nil {
			continue
		}

		// -------------------
		// Filter logs
		// -------------------
		if level != "" && entry.Level != level {
			continue
		}
		if method != "" && entry.Method != method {
			continue
		}
		if path != "" && entry.Path != path {
			continue
		}
		if status != "" {
			if stInt, _ := strconv.Atoi(status); entry.Status != stInt {
				continue
			}
		}

		logs = append(logs, entry)
	}

	// -------------------
	// Pagination
	// -------------------
	total := len(logs)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	if page < 1 {
		page = 1
	}

	// vị trí bắt đầu: từ cuối mảng đếm ngược
	start := max(total-page*limit, 0)

	end := min(start+limit, total)

	pageLogs := logs[start:end]

	// nhưng bạn muốn newest -> oldest
	// nên đảo nhỏ gọn subset (chỉ đảo ~50 dòng)
	for i, j := 0, len(pageLogs)-1; i < j; i, j = i+1, j-1 {
		pageLogs[i], pageLogs[j] = pageLogs[j], pageLogs[i]
	}

	c.JSON(200, gin.H{
		"page":  page,
		"limit": limit,
		"total": total,
		"data":  pageLogs,
	})
}
