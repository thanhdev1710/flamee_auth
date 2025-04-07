package services

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhdev1710/flamee_auth/internal/models"
	"github.com/thanhdev1710/flamee_auth/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

type UserServices struct {
	userRepo              *repo.UserRepo
	verificationTokenRepo *repo.VerificationTokenRepo
	emailServices         *EmailServices
	sessionRepo           *repo.SessionRepo
}

func NewUserServices() *UserServices {
	return &UserServices{
		userRepo:              repo.NewUserRepo(),
		verificationTokenRepo: repo.NewVerificationTokenRepo(),
		emailServices:         NewEmailServices(),
		sessionRepo:           repo.NewSessionRepo(),
	}
}

func (us *UserServices) ConfirmEmail(email string) error {
	user, err := us.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}

	if user.IsVerified {
		return errors.New("email already verified")
	}

	user.IsVerified = true

	err = us.userRepo.Save(user)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserServices) SendResetPassword(email string, c *gin.Context) error {
	user, err := repo.NewUserRepo().FindByEmail(email)
	if err != nil {
		return err
	}

	token := uuid.New().String()

	verificationToken := models.VerificationToken{
		UserId:    user.Id,
		Token:     token,
		TokenType: "password_reset",
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	err = us.verificationTokenRepo.Create(&verificationToken)
	if err != nil {
		return err
	}

	// Tạo URL xác thực chứa token
	protocol := "http"
	if c.Request.TLS != nil {
		protocol = "https"
	}

	verificationURL := fmt.Sprintf("%s://%s/auth/change-password/%s", protocol, c.Request.Host, token)

	// Gửi email xác nhận
	us.emailServices.Send(email, verificationURL, "password_reset")
	return nil
}

func (us *UserServices) UpdatePassword(token string, password string) error {
	// Tìm token trong repository
	vToken, err := us.verificationTokenRepo.FindByToken(token, "password_reset")
	if err != nil {
		return errors.New("invalid token")
	}

	// Kiểm tra thời gian hết hạn của token
	if time.Now().After(vToken.ExpiresAt) {
		return errors.New("token has expired")
	}

	// Mã hóa mật khẩu mới
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(3)

	// Tạo channel lỗi để nhận lỗi từ các goroutines
	errCh := make(chan error, 3) // Có thể chứa nhiều lỗi từ các goroutines

	userId := vToken.UserId

	// Cập nhật mật khẩu của người dùng trong cơ sở dữ liệu
	go func() {
		defer wg.Done()
		err := us.userRepo.UpdatePassword(userId, string(hashedPassword))
		if err != nil {
			errCh <- fmt.Errorf("failed to update password: %w", err)
		}
	}()

	// Xóa verification token vì đã sử dụng xong
	go func() {
		defer wg.Done()
		err := us.verificationTokenRepo.Delete(vToken)
		if err != nil {
			errCh <- fmt.Errorf("failed to delete verification token: %w", err)
		}
	}()

	// Thu hồi tất cả session của người dùng
	go func() {
		defer wg.Done()
		err := us.sessionRepo.RevokeTokensByUserId(userId)
		if err != nil {
			errCh <- fmt.Errorf("failed to revoke sessions: %w", err)
		}
	}()

	// Chờ các goroutines hoàn thành
	wg.Wait()
	close(errCh)

	// Kiểm tra lỗi trong channel và trả về lỗi đầu tiên
	for err := range errCh {
		if err != nil {
			return err // Trả về lỗi đầu tiên gặp phải
		}
	}

	return nil
}
