package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

	// Gọi AuthService để xử lý logic refresh token
	accessToken, err := ac.authServices.RefreshToken(cookieToken, claims, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// Trả về access token mới
	c.JSON(http.StatusOK, gin.H{
		"message": "Refresh token success",
		"token":   accessToken,
	})
}

func (ac *AuthControllers) Logout(c *gin.Context) {
	// Gọi service LogoutUser để thực hiện đăng xuất
	err := ac.authServices.LogoutUser(c)
	if err != nil {
		// Nếu có lỗi xảy ra, trả về lỗi với thông báo và mã lỗi phù hợp
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Logout failed: " + err.Error(),
		})
		return
	}

	// Nếu đăng xuất thành công, trả về phản hồi thành công
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}
