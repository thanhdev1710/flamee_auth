package repo

import (
	"errors"
	"fmt"

	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/internal/models"
)

type VerificationTokenRepo struct{}

func NewVerificationTokenRepo() *VerificationTokenRepo {
	return &VerificationTokenRepo{}
}

func (vr *VerificationTokenRepo) Create(verificationToken *models.VerificationToken) error {
	return global.Pdb.Create(verificationToken).Error
}

func (vr *VerificationTokenRepo) Save(verificationToken *models.VerificationToken) error {
	return global.Pdb.Save(verificationToken).Error
}

func (vr *VerificationTokenRepo) Delete(verificationToken *models.VerificationToken) error {
	return global.Pdb.Delete(verificationToken).Error
}

func (vr *VerificationTokenRepo) FindByToken(token string, tokenType string) (*models.VerificationToken, error) {
	var verificationToken models.VerificationToken
	// Chỉnh lại điều kiện WHERE cho đúng với cột tokenType thay vì 'password_reset'
	fmt.Println(token, tokenType)
	if err := global.Pdb.Where("token_type = ? AND token = ?", tokenType, token).First(&verificationToken).Error; err != nil {
		return nil, errors.New("verification token not found")
	}

	return &verificationToken, nil
}
