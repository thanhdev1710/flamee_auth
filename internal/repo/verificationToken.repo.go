package repo

import (
	"errors"

	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/internal/models"
	"go.uber.org/zap"
)

type VerificationTokenRepo struct{}

func NewVerificationTokenRepo() *VerificationTokenRepo {
	return &VerificationTokenRepo{}
}

func (vr *VerificationTokenRepo) Create(verificationToken *models.VerificationToken) error {
	if err := global.Pdb.Create(verificationToken).Error; err != nil {
		global.Logger.Error("Lỗi khi tạo verification token", zap.Error(err))
		return errors.New("lỗi máy chủ, không thể tạo mã xác thực")
	}
	return nil
}

func (vr *VerificationTokenRepo) Save(verificationToken *models.VerificationToken) error {
	if err := global.Pdb.Save(verificationToken).Error; err != nil {
		global.Logger.Error("Lỗi khi lưu verification token", zap.Error(err))
		return errors.New("lỗi máy chủ, không thể lưu mã xác thực")
	}
	return nil
}

func (vr *VerificationTokenRepo) Delete(verificationToken *models.VerificationToken) error {
	if err := global.Pdb.Delete(verificationToken).Error; err != nil {
		global.Logger.Error("Lỗi khi xóa verification token", zap.Error(err))
		return errors.New("lỗi máy chủ, không thể xóa mã xác thực")
	}
	return nil
}

func (vr *VerificationTokenRepo) FindByToken(token string, tokenType string) (*models.VerificationToken, error) {
	var verificationToken models.VerificationToken
	if err := global.Pdb.Where("token_type = ? AND token = ?", tokenType, token).First(&verificationToken).Error; err != nil {
		if errors.Is(err, global.Pdb.Error) || err.Error() == "record not found" {
			return nil, errors.New("mã xác thực không hợp lệ hoặc không tồn tại")
		}
		global.Logger.Error("Lỗi khi tìm verification token", zap.Error(err), zap.String("token", token), zap.String("token_type", tokenType))
		return nil, errors.New("lỗi máy chủ, vui lòng thử lại sau")
	}
	return &verificationToken, nil
}
