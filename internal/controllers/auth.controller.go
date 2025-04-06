package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/internal/models"
	"github.com/thanhdev1710/flamee_auth/internal/services"
	"github.com/thanhdev1710/flamee_auth/pkg/utils"
)

type AuthControllers struct {
	authServices *services.AuthServices
}

func NewAuthControllers() *AuthControllers {
	return &AuthControllers{
		authServices: services.NewAuthServices(),
	}
}

func (ac *AuthControllers) Register(c *gin.Context) {
	var user services.UserRegisterRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	token, err := ac.authServices.RegisterUser(user, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"token":   token,
	})
}

func (ac *AuthControllers) Login(c *gin.Context) {
	var user services.UserLoginRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Đăng nhập và tạo token
	token, err := ac.authServices.LoginUser(user, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	// Trả về token nếu đăng nhập thành công
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

func (ac *AuthControllers) RefreshToken(c *gin.Context) {
	// Lấy refresh token từ cookie
	cookieToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Refresh token missing"})
		return
	}

	// Xác thực refresh token
	claims, err := utils.ValidateToken(cookieToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid refresh token"})
		return
	}

	// Kiểm tra trong DB session
	var session models.Session
	if err := global.Pdb.
		Where("user_id = ? AND token = ?", claims.Subject, cookieToken).
		First(&session).Error; err != nil {
		// Nếu không tìm thấy session hoặc session hết hạn
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Refresh token expired or invalid"})
		return
	}

	// Kiểm tra xem refresh token có hết hạn không
	if session.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Refresh token expired"})
		return
	}

	// Lấy thông tin người dùng từ database
	var user models.User
	if err := global.Pdb.First(&user, "id = ?", claims.Subject).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
		return
	}

	// Tạo access token mới
	tokenDuration, err := time.ParseDuration(global.Config.JwtExpirationTimeDefault)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse token expiration time"})
		return
	}
	accessToken, err := utils.GenerateToken(&user, tokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate access token"})
		return
	}

	// Tạo refresh token mới nếu muốn "rotate" refresh token
	newRefreshDuration, err := time.ParseDuration(global.Config.JwtExpirationTimeRemember)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse refresh token expiration time"})
		return
	}
	newRefreshToken, err := utils.GenerateToken(&user, newRefreshDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate refresh token"})
		return
	}

	// Cập nhật refresh token trong session nếu cần
	session.Token = newRefreshToken
	session.ExpiresAt = time.Now().Add(newRefreshDuration)
	if err := global.Pdb.Save(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update session"})
		return
	}

	utils.SetCookiesToken(c, accessToken, newRefreshToken, tokenDuration, newRefreshDuration)
	// Trả về access token mới
	c.JSON(http.StatusOK, gin.H{
		"message": "Refresh token success",
		"token":   accessToken,
	})
}
