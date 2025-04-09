package utils

import "regexp"

func IsValidEmail(email string) bool {
	// Biểu thức chính quy để kiểm tra email hợp lệ
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func IsValidPassword(password string) bool {
	// Kiểm tra độ dài tối thiểu
	if len(password) < 13 {
		return false
	}

	// Biểu thức chính quy kiểm tra từng yêu cầu
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]{};':"\\|,.<>\/?]+`).MatchString(password)

	return hasLower && hasUpper && hasNumber && hasSpecial
}
