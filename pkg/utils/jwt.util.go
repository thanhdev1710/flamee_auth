package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/internal/models"
)

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// Secret key để ký JWT (trong thực tế bạn nên lưu trữ trong biến môi trường hoặc cấu hình)
var jwtSecret = []byte(global.Config.JwtSecret)

// Tạo JWT token mới
func GenerateToken(user *models.User, parsedExpirationTime time.Duration) (string, error) {
	claims := &Claims{
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Id.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(parsedExpirationTime)), // Token hết hạn sau 24 giờ
			Issuer:    "flamee_auth",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Tạo token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// Xác thực token
func ValidateToken(tokenStr string) (*Claims, error) {
	// Parse và xác thực token
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
		// Kiểm tra kiểu signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	// Trả về claims nếu token hợp lệ
	if claims, ok := token.Claims.(*Claims); ok {
		// Kiểm tra thời gian hết hạn của token
		if claims.ExpiresAt.Before(time.Now()) {
			return nil, errors.New("token has expired")
		}
		return claims, nil
	}
	return nil, errors.New("invalid claims")
}
