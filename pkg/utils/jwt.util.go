package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/internal/models"
)

type Claims struct {
	Role  string `json:"role"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// Tạo JWT token mới
func GenerateToken(user *models.User, parsedExpirationTime time.Duration) (string, error) {
	claims := &Claims{
		Role:  user.Role,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Id.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(parsedExpirationTime)), // Token hết hạn sau 24 giờ
			Issuer:    "flamee_auth",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Tạo token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err := token.SignedString([]byte(global.Config.JwtSecret))

	if err != nil {
		return "", errors.New("lỗi khi tạo token")
	}

	return signedToken, nil
}

func ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("phương thức ký không hợp lệ")
		}
		return []byte(global.Config.JwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("token không hợp lệ hoặc bị lỗi")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("thông tin trong token không hợp lệ")
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token đã hết hạn")
	}

	return claims, nil
}

func ParseDuration(timeStr string) (time.Duration, error) {
	timeValue, err := time.ParseDuration(timeStr)
	if err != nil {
		return 0, errors.New("cấu hình thời gian không hợp lệ")
	}

	return timeValue, nil
}

func UuidParse(uuidStr string) (uuid.UUID, error) {

	userId, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.Nil, errors.New("định dạng id người dùng không hợp lệ")
	}

	return userId, nil
}
