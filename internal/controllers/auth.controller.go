package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/internal/repo"
	"github.com/thanhdev1710/flamee_auth/internal/services"
	"github.com/thanhdev1710/flamee_auth/pkg/utils"
)

type AuthControllers struct {
	authServices  *services.AuthServices
	emailServices *services.EmailServices
	userServices  *services.UserServices
}

func NewAuthControllers() *AuthControllers {
	return &AuthControllers{
		authServices:  services.NewAuthServices(),
		emailServices: services.NewEmailServices(),
		userServices:  services.NewUserServices(),
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

func (ac *AuthControllers) SendVerifyEmail(c *gin.Context) {
	email := c.Param("email")
	// Kiểm tra xem email có hợp lệ không
	if email == "" || !utils.IsValidEmail(email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid or missing email parameter",
		})
		return
	}

	user, err := repo.NewUserRepo().FindByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if user.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email already verified",
		})
		return
	}

	token, err := utils.Encrypt(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate verification token",
		})
		return
	}

	// Tạo URL xác thực chứa token
	protocol := "http"
	if c.Request.TLS != nil {
		protocol = "https"
	}

	verificationURL := fmt.Sprintf("%s://%s/auth/verify-email/%s", protocol, c.Request.Host, token)

	// Gửi email xác nhận
	ac.emailServices.SendVerificationEmail(email, verificationURL)

	// Phản hồi về việc gửi email
	c.JSON(http.StatusOK, gin.H{
		"message": "Verification email is being sent. Please check your inbox within 24 hours.",
	})
}

func (ac *AuthControllers) VerifyEmail(c *gin.Context) {
	// Lấy token từ tham số trong URL
	token := c.Param("token")

	// Giải mã token để lấy email
	email, err := utils.Decrypt(token)
	if err != nil {
		// Nếu lỗi giải mã token, trả về lỗi với thông báo tương ứng
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token không hợp lệ hoặc đã hết hạn",
		})
		return
	}

	// Gọi ConfirmEmail để cập nhật trạng thái xác thực email
	err = ac.userServices.ConfirmEmail(email)
	if err != nil {
		// Nếu có lỗi khi xác thực email, trả về lỗi với thông báo tương ứng
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Nếu email được xác thực thành công, trả về phản hồi thành công
	c.JSON(http.StatusOK, gin.H{
		"email":   email,
		"message": "Email đã được xác thực thành công",
	})
}
