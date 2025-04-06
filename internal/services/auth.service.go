package services

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/internal/models"
	"github.com/thanhdev1710/flamee_auth/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRegisterRequest là cấu trúc nhận dữ liệu từ client khi đăng ký
type UserRegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// UserLoginRequest là cấu trúc nhận dữ liệu từ client khi đăng nhập
type UserLoginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe"`
}

type AuthServices struct{}

func NewAuthServices() *AuthServices {
	return &AuthServices{}
}

func (as *AuthServices) RegisterUser(user UserRegisterRequest, ctx *gin.Context) (string, error) {
	var existingUser models.User
	result := global.Pdb.Where("email = ?", user.Email).First(&existingUser)

	if result.Error == nil {
		return "", errors.New("user already exists")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", errors.New("internal server error")
	}

	// Mã hóa mật khẩu
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	newUser := models.User{
		Id:       uuid.New(),
		Email:    user.Email,
		Password: string(hashedPassword),
		Role:     user.Role,
	}

	if err := global.Pdb.Create(&newUser).Error; err != nil {
		return "", err
	}

	// Tạo access token
	timeDefault, err := time.ParseDuration(global.Config.JwtExpirationTimeDefault)
	if err != nil {
		return "", err
	}

	accessToken, err := utils.GenerateToken(&newUser, timeDefault)
	if err != nil {
		return "", err
	}

	// Auto remember (nếu muốn) sau khi đăng ký
	timeRemember, _ := time.ParseDuration(global.Config.JwtExpirationTimeRemember)
	refreshToken, err := utils.GenerateToken(&newUser, timeRemember)
	if err != nil {
		return "", err
	}

	session := models.Session{
		Token:     refreshToken,
		UserId:    newUser.Id,
		UserAgent: ctx.Request.UserAgent(),
		IpAddress: ctx.ClientIP(),
		ExpiresAt: time.Now().Add(timeRemember),
	}

	if err := global.Pdb.Create(&session).Error; err != nil {
		return "", err
	}

	// Set token cookie
	utils.SetCookiesToken(ctx, accessToken, refreshToken, timeDefault, timeRemember)
	return accessToken, nil
}

func (as *AuthServices) LoginUser(user UserLoginRequest, ctx *gin.Context) (string, error) {
	// Lấy người dùng từ cơ sở dữ liệu
	userFromDB := models.User{}
	if err := global.Pdb.Where("email = ?", user.Email).First(&userFromDB).Error; err != nil {
		return "", errors.New("user not found")
	}

	// Kiểm tra mật khẩu
	if err := bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(user.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Tạo access token
	timeDefault, _ := time.ParseDuration(global.Config.JwtExpirationTimeDefault)
	accessToken, err := utils.GenerateToken(&userFromDB, timeDefault)
	if err != nil {
		return "", err
	}

	// Nếu có Remember Me thì tạo thêm refresh token và lưu session
	if user.RememberMe {
		timeRemember, _ := time.ParseDuration(global.Config.JwtExpirationTimeRemember)
		refreshToken, err := utils.GenerateToken(&userFromDB, timeRemember)
		if err != nil {
			return "", err
		}

		session := models.Session{
			Token:     refreshToken,
			UserId:    userFromDB.Id,
			UserAgent: ctx.Request.UserAgent(),
			IpAddress: ctx.ClientIP(),
			ExpiresAt: time.Now().Add(timeRemember),
		}

		if err := global.Pdb.Create(&session).Error; err != nil {
			return "", err
		}

		// Set token cookie
		utils.SetCookiesToken(ctx, accessToken, refreshToken, timeDefault, timeRemember)
	}

	return accessToken, nil
}
