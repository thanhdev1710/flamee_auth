package utils

import "regexp"

func IsValidEmail(email string) bool {
	// Biểu thức chính quy để kiểm tra email hợp lệ
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
