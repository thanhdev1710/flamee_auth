package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/global"
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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "một số trường chưa hợp lệ, vui lòng kiểm tra lại",
		})
		return
	}

	if !utils.IsValidEmail(user.Email) || !utils.IsValidPassword(user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "một số trường chưa hợp lệ, vui lòng kiểm tra lại",
		})
		return
	}

	token, err := ac.authServices.RegisterUser(user, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Tạo URL xác thực chứa token
	verificationURL := fmt.Sprintf("%s/auth/verify-email/%s", global.Url.UrlFrontEnd, token)
	// Gửi email xác nhận
	ac.emailServices.Send(user.Email, verificationURL, "verification")

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Đăng ký tài khoản thành công",
		"token":   token,
	})
}

func (ac *AuthControllers) Login(c *gin.Context) {
	var user services.UserLoginRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "một số trường chưa hợp lệ, vui lòng kiểm tra lại",
		})
		return
	}

	if !utils.IsValidEmail(user.Email) || !utils.IsValidPassword(user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "một số trường chưa hợp lệ, vui lòng kiểm tra lại",
		})
		return
	}

	// Đăng nhập và tạo token
	token, err := ac.authServices.LoginUser(user, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Trả về token nếu đăng nhập thành công
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Đăng nhập thành công",
		"token":   token,
	})
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (ac *AuthControllers) RefreshToken(c *gin.Context) {
	// Lấy refresh token từ cookie
	refreshToken, err := c.Cookie(utils.HexString(global.Token.RefreshToken))
	if err != nil || refreshToken == "" {
		var req RefreshTokenRequest
		if bindErr := c.ShouldBindJSON(&req); bindErr != nil || req.RefreshToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Không tìm thấy refresh token trong cookie",
			})
			return
		}

		refreshToken = req.RefreshToken
	}

	// Xác thực refresh token
	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Gọi AuthService để xử lý logic refresh token
	accessToken, err := ac.authServices.RefreshToken(refreshToken, claims, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Trả về access token mới
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Làm mới token thành công",
		"token":   accessToken,
	})
}

func (ac *AuthControllers) Logout(c *gin.Context) {
	// Gọi service LogoutUser để thực hiện đăng xuất
	err := ac.authServices.LogoutUser(c)
	if err != nil {
		// Nếu có lỗi xảy ra, trả về lỗi với thông báo và mã lỗi phù hợp
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Nếu đăng xuất thành công, trả về phản hồi thành công
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Đăng xuất thành công",
	})
}

func (ac *AuthControllers) SendVerifyEmail(c *gin.Context) {
	email := c.Param("email")
	// Kiểm tra xem email có hợp lệ không
	if email == "" || !utils.IsValidEmail(email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "một số trường chưa hợp lệ, vui lòng kiểm tra lại",
		})
		return
	}

	user, err := repo.NewUserRepo().FindByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	if user.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Email này đã xác thực",
		})
		return
	}

	token, err := utils.Encrypt(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Tạo URL xác thực chứa token

	verificationURL := fmt.Sprintf("%s/auth/verify-email/%s", global.Url.UrlFrontEnd, token)

	// Gửi email xác nhận
	ac.emailServices.Send(email, verificationURL, "verification")

	// Phản hồi về việc gửi email
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Email xác minh đang được gửi. Vui lòng kiểm tra hộp thư đến của bạn trong vòng 24 giờ.",
	})
}

func (ac *AuthControllers) VerifyEmail(c *gin.Context) {
	// Lấy token từ tham số trong URL
	token := c.Param("token")

	fmt.Println("Token ::", token)
	// Giải mã token để lấy email
	email, err := utils.Decrypt(token)
	if err != nil {
		// Nếu lỗi giải mã token, trả về lỗi với thông báo tương ứng
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Gọi ConfirmEmail để cập nhật trạng thái xác thực email
	err = ac.userServices.ConfirmEmail(email)
	if err != nil {
		// Nếu có lỗi khi xác thực email, trả về lỗi với thông báo tương ứng
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Nếu email được xác thực thành công, trả về phản hồi thành công
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Xác thực tài khoản thành công",
		"email":   email,
	})
}

func (ac *AuthControllers) SendResetPassword(c *gin.Context) {
	email := c.Param("email")
	// Kiểm tra xem email có hợp lệ không
	if email == "" || !utils.IsValidEmail(email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "một số trường chưa hợp lệ, vui lòng kiểm tra lại",
		})
		return
	}

	err := ac.userServices.SendResetPassword(email, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Email đặt lại mật khẩu đang được gửi. Vui lòng kiểm tra hộp thư đến của bạn trong vòng 5 phút.",
	})
}

type ResetPasswordRequest struct {
	Password string `json:"password"`
	Token    string `json:"token"`
}

func (ac *AuthControllers) ResetPassword(c *gin.Context) {
	var body ResetPasswordRequest
	if err := c.ShouldBindJSON(&body); err != nil || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "một số trường chưa hợp lệ, vui lòng kiểm tra lại",
		})
		return
	}

	if !utils.IsValidPassword(body.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "một số trường chưa hợp lệ, vui lòng kiểm tra lại",
		})
		return
	}

	err := ac.userServices.UpdatePassword(body.Token, body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Đặt lại mật khẩu thành công",
	})
}

type DeleteAccountRequest struct {
	Password string `json:"password"`
}

func (ac *AuthControllers) DeleteAccount(c *gin.Context) {
	userId := utils.GetUserId(c)
	var body DeleteAccountRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "một số trường chưa hợp lệ, vui lòng kiểm tra lại",
		})
		return
	}

	if err := ac.userServices.DeleteAccount(userId, body.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Xoá tài khoản thành công",
	})
}

func (ac *AuthControllers) RestoreAccount(c *gin.Context) {
	userId := utils.GetUserId(c)

	if err := ac.userServices.RestoreAccount(userId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Khôi phục tài khoản thành công",
	})
}
