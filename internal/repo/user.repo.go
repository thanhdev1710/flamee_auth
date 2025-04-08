package repo

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/internal/models"
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
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepo) FindById(id string) (*models.User, error) {
	var user models.User

	err := global.Pdb.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepo) Create(user *models.User) error {
	return global.Pdb.Create(user).Error
}

func (ur *UserRepo) Save(user *models.User) error {
	return global.Pdb.Save(user).Error
}

func (ur *UserRepo) Updates(user *models.User, value any) error {
	return global.Pdb.Model(user).Updates(value).Error
}

func (ur *UserRepo) UpdatePassword(userId uuid.UUID, password string) error {
	if err := global.Pdb.Model(&models.User{}).Where("id = ?", userId).Updates(map[string]any{
		"password":   password,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return err
	}
	return nil
}
