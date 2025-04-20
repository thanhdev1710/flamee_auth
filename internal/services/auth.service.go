package services

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/internal/models"
	"github.com/thanhdev1710/flamee_auth/internal/repo"
	"github.com/thanhdev1710/flamee_auth/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

// UserRegisterRequest là cấu trúc nhận dữ liệu từ client khi đăng ký
type UserRegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserLoginRequest là cấu trúc nhận dữ liệu từ client khi đăng nhập
type UserLoginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe"`
}

type AuthServices struct {
	userRepo    *repo.UserRepo
	sessionRepo *repo.SessionRepo
}

func NewAuthServices() *AuthServices {
	return &AuthServices{
		userRepo:    repo.NewUserRepo(),
		sessionRepo: repo.NewSessionRepo(),
	}
}

func (as *AuthServices) RegisterUser(user UserRegisterRequest, c *gin.Context) (string, error) {
	_, err := as.userRepo.FindByEmail(user.Email)
	if err == nil {
		return "", errors.New("tài khoản này đã tồn tại")
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
	}

	if err := as.userRepo.Create(&newUser); err != nil {
		return "", err
	}

	timeDefault, err := utils.ParseDuration(global.Config.JwtExpirationTimeDefault)
	if err != nil {
		return "", err
	}
	timeRemember, err := time.ParseDuration(global.Config.JwtExpirationTimeRemember)
	if err != nil {
		return "", err
	}

	accessToken, err := utils.GenerateToken(&newUser, timeDefault)
	if err != nil {
		return "", err
	}
	refreshToken, err := utils.GenerateToken(&newUser, timeRemember)
	if err != nil {
		return "", err
	}

	session := models.Session{
		Token:     refreshToken,
		UserId:    newUser.Id,
		UserAgent: c.Request.UserAgent(),
		IpAddress: c.ClientIP(),
		ExpiresAt: time.Now().Add(timeRemember),
	}

	if err := as.sessionRepo.Create(&session); err != nil {
		return "", err
	}

	// Set token cookie
	utils.SetCookiesToken(c, accessToken, refreshToken, timeDefault, timeRemember)
	return accessToken, nil
}

func (as *AuthServices) LoginUser(user UserLoginRequest, c *gin.Context) (string, error) {
	// Lấy người dùng từ cơ sở dữ liệu
	userFromDB, err := as.userRepo.FindByEmail(user.Email)
	if err != nil {
		return "", err
	}

	// Kiểm tra mật khẩu
	if err := utils.CompareHashAndPassword(userFromDB.Password, user.Password); err != nil {
		return "", err
	}

	// Tạo access token
	timeDefault, err := utils.ParseDuration(global.Config.JwtExpirationTimeDefault)
	if err != nil {
		return "", err
	}

	accessToken, err := utils.GenerateToken(userFromDB, timeDefault)
	if err != nil {
		return "", err
	}

	// Nếu có Remember Me thì tạo thêm refresh token và lưu session
	var timeRememberTmp time.Duration
	var refreshTokenTmp string
	if user.RememberMe {
		timeRemember, err := utils.ParseDuration(global.Config.JwtExpirationTimeRemember)
		if err != nil {
			return "", err
		}

		refreshToken, err := utils.GenerateToken(userFromDB, timeRemember)
		if err != nil {
			return "", err
		}

		session := models.Session{
			Token:     refreshToken,
			UserId:    userFromDB.Id,
			UserAgent: c.Request.UserAgent(),
			IpAddress: c.ClientIP(),
			ExpiresAt: time.Now().Add(timeRemember),
		}

		if err := as.sessionRepo.Create(&session); err != nil {
			return "", err
		}

		timeRememberTmp = timeRemember
		refreshTokenTmp = refreshToken
	}

	// Set token cookie
	utils.SetCookiesToken(c, accessToken, refreshTokenTmp, timeDefault, timeRememberTmp)
	return accessToken, nil
}

func (as *AuthServices) RefreshToken(cookieToken string, claims *utils.Claims, c *gin.Context) (string, error) {
	// Kiểm tra trong DB session
	session, err := as.sessionRepo.FindByUserAndToken(claims.Subject, cookieToken)
	if err != nil {
		return "", err
	}

	// Lấy thông tin người dùng từ database
	user, err := as.userRepo.FindById(claims.Subject)
	if err != nil {
		return "", err
	}

	// Tạo access token mới
	timeDefault, err := utils.ParseDuration(global.Config.JwtExpirationTimeDefault)
	if err != nil {
		return "", err
	}
	accessToken, err := utils.GenerateToken(user, timeDefault)
	if err != nil {
		return "", err
	}

	// Tạo refresh token mới nếu muốn "rotate" refresh token
	timeRemember, err := utils.ParseDuration(global.Config.JwtExpirationTimeRemember)
	if err != nil {
		return "", err
	}
	newRefreshToken, err := utils.GenerateToken(user, timeRemember)
	if err != nil {
		return "", err
	}

	// Cập nhật refresh token trong session nếu cần
	session.Token = newRefreshToken
	session.ExpiresAt = time.Now().Add(timeRemember)

	if err := as.sessionRepo.Save(session); err != nil {
		return "", err
	}

	utils.SetCookiesToken(c, accessToken, newRefreshToken, timeDefault, timeRemember)
	return accessToken, nil
}

func (as *AuthServices) LogoutUser(c *gin.Context) error {
	userIdStr := utils.GetUserId(c)
	// Parse string thành UUID
	userId, err := utils.UuidParse(userIdStr)
	if err != nil {
		return err
	}

	utils.SetCookiesToken(c, "", "", -1, -1)

	return as.sessionRepo.RevokeTokensByUserId(userId)
}
