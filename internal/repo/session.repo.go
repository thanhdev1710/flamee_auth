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

type SessionRepo struct{}

func NewSessionRepo() *SessionRepo {
	return &SessionRepo{}
}

func (sr *SessionRepo) Create(session *models.Session) error {
	if err := global.Pdb.Create(session).Error; err != nil {
		global.Logger.Error("Lỗi khi tạo session", zap.Error(err))
		return errors.New("lỗi máy chủ, vui lòng thử lại sau")
	}
	return nil
}

func (sr *SessionRepo) Save(session *models.Session) error {
	if err := global.Pdb.Save(session).Error; err != nil {
		global.Logger.Error("Lỗi khi lưu session", zap.Error(err))
		return errors.New("lỗi máy chủ, vui lòng thử lại sau")
	}
	return nil
}

func (sr *SessionRepo) FindByUserAndToken(userId, token string) (*models.Session, error) {
	var session models.Session
	err := global.Pdb.Where("user_id = ? AND token = ?", userId, token).First(&session).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("refresh token không hợp lệ")
	} else if err != nil {
		global.Logger.Error("Lỗi khi tìm session theo userId và token", zap.Error(err))
		return nil, errors.New("lỗi máy chủ, vui lòng thử lại sau")
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token đã hết hạn")
	}

	return &session, nil
}

func (sr *SessionRepo) RevokeTokensByUserId(userId uuid.UUID) error {
	if err := global.Pdb.Where("user_id = ?", userId).Delete(&models.Session{}).Error; err != nil {
		global.Logger.Error("Lỗi khi thu hồi tất cả phiên đăng nhập của người dùng", zap.Error(err), zap.String("user_id", userId.String()))
		return errors.New("lỗi máy chủ, không thể thu hồi token")
	}
	return nil
}
