package repo

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepo struct {
}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (ur *UserRepo) FindByEmail(email string) (*models.User, error) {
	var user models.User

	err := global.Pdb.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("tài khoản không tồn tại")
	} else if err != nil {
		global.Logger.Error("Lỗi khi tìm người dùng theo email", zap.Error(err), zap.String("email", email))
		return nil, errors.New("lỗi máy chủ, vui lòng thử lại sau")
	}

	return &user, nil
}

func (ur *UserRepo) FindById(id string) (*models.User, error) {
	var user models.User

	err := global.Pdb.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("tài khoản không tồn tại")
	} else if err != nil {
		global.Logger.Error("Lỗi khi tìm người dùng theo ID", zap.Error(err), zap.String("id", id))
		return nil, errors.New("lỗi máy chủ, vui lòng thử lại sau")
	}

	return &user, nil
}

func (ur *UserRepo) Create(user *models.User) error {
	if err := global.Pdb.Create(user).Error; err != nil {
		global.Logger.Error("Lỗi khi tạo người dùng", zap.Error(err), zap.String("email", user.Email))
		return errors.New("lỗi máy chủ, không thể tạo người dùng")
	}
	return nil
}

func (ur *UserRepo) Save(user *models.User) error {
	if err := global.Pdb.Save(user).Error; err != nil {
		global.Logger.Error("Lỗi khi lưu người dùng", zap.Error(err), zap.String("id", user.Id.String()))
		return errors.New("lỗi máy chủ, không thể lưu người dùng")
	}
	return nil
}

func (ur *UserRepo) Updates(user *models.User, value any) error {
	if err := global.Pdb.Model(user).Updates(value).Error; err != nil {
		global.Logger.Error("Lỗi khi cập nhật người dùng", zap.Error(err), zap.String("id", user.Id.String()))
		return errors.New("lỗi máy chủ, không thể cập nhật người dùng")
	}
	return nil
}

func (ur *UserRepo) UpdatePassword(userId uuid.UUID, password string) error {
	if err := global.Pdb.Model(&models.User{}).Where("id = ?", userId).Updates(map[string]any{
		"password":   password,
		"updated_at": time.Now(),
	}).Error; err != nil {
		global.Logger.Error("Lỗi khi cập nhật mật khẩu", zap.Error(err), zap.String("user_id", userId.String()))
		return errors.New("lỗi máy chủ, không thể cập nhật mật khẩu")
	}
	return nil
}
