package services

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/internal/models"
	"github.com/thanhdev1710/flamee_auth/internal/repo"
	"github.com/thanhdev1710/flamee_auth/pkg/utils"
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
		return errors.New("tài khoản này đã xác thực")
	}

	user.IsVerified = true
	user.Status = global.User.Active
	user.UpdatedAt = time.Now()

	err = us.userRepo.Save(user)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserServices) ConfirmProfile(userId string) error {
	user, err := us.userRepo.FindById(userId)
	if err != nil {
		return err
	}

	if user.IsProfile {
		return errors.New("tài khoản này đã có")
	}

	user.IsProfile = true
	user.UpdatedAt = time.Now()

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

	verificationToken := models.VerificationToken{
		UserId:    user.Id,
		Token:     uuid.New().String(),
		TokenType: "password_reset",
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	err = us.verificationTokenRepo.Create(&verificationToken)
	if err != nil {
		return err
	}

	// Tạo URL xác thực chứa token
	verificationURL := fmt.Sprintf("%s/auth/reset-password/%s", global.Url.UrlFrontEnd, verificationToken.Token)
	fmt.Println(verificationURL)
	// Gửi email xác nhận
	us.emailServices.Send(email, verificationURL, "password_reset")
	return nil
}

func (us *UserServices) UpdatePassword(token string, password string) error {
	// Tìm token trong repository
	vToken, err := us.verificationTokenRepo.FindByToken(token, "password_reset")
	if err != nil {
		return err
	}

	// Kiểm tra thời gian hết hạn của token
	if time.Now().After(vToken.ExpiresAt) {
		return errors.New("token đã hết hạn")
	}

	// Mã hóa mật khẩu mới
	hashedPassword, err := utils.GenerateFromPassword(password)
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
			errCh <- err
		}
	}()

	// Xóa verification token vì đã sử dụng xong
	go func() {
		defer wg.Done()
		err := us.verificationTokenRepo.Delete(vToken)
		if err != nil {
			errCh <- err
		}
	}()

	// Thu hồi tất cả session của người dùng
	go func() {
		defer wg.Done()
		err := us.sessionRepo.RevokeTokensByUserId(userId)
		if err != nil {
			errCh <- err
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

func (us *UserServices) DeleteAccount(userId string, password string) error {
	// Tìm user theo ID
	existingUser, err := us.userRepo.FindById(userId)
	if err != nil {
		return err
	}

	if err := utils.CompareHashAndPassword(existingUser.Password, password); err != nil {
		return err
	}

	if existingUser.Status != global.User.Active {
		return fmt.Errorf("không thể khôi phục tài khoản với trạng thái: %s", existingUser.Status)
	}

	// Cập nhật trạng thái user thành "banned"
	now := time.Now()
	updateData := models.User{
		Status:    global.User.Banned,
		UpdatedAt: now,
		DeletedAt: &now,
	}

	if err := us.userRepo.Updates(existingUser, updateData); err != nil {
		return fmt.Errorf("xoá tài khoản thất bại: %w", err)
	}

	return nil
}

func (us *UserServices) RestoreAccount(userId string) error {
	// Tìm user theo ID
	existingUser, err := us.userRepo.FindById(userId)
	if err != nil {
		return err
	}

	if existingUser.Status != global.User.Banned {
		return fmt.Errorf("không thể khôi phục tài khoản với trạng thái: %s", existingUser.Status)
	}

	// Cập nhật trạng thái user thành "active"
	updateData := models.User{
		Status:    global.User.Active,
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}

	if err := us.userRepo.Updates(existingUser, updateData); err != nil {
		return err
	}

	return nil
}
