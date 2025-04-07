package repo

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/internal/models"
	"gorm.io/gorm"
)

type SessionRepo struct{}

func NewSessionRepo() *SessionRepo {
	return &SessionRepo{}
}

func (sr *SessionRepo) Create(session *models.Session) error {
	return global.Pdb.Create(session).Error
}

func (sr *SessionRepo) Save(session *models.Session) error {
	return global.Pdb.Save(session).Error
}

func (sr *SessionRepo) FindByUserAndToken(userId, token string) (*models.Session, error) {
	var session models.Session
	err := global.Pdb.Where("user_id = ? AND token = ?").First(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("refresh token invalid")
	} else if err != nil {
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token expired")
	}

	return &session, nil
}

func (sr *SessionRepo) RevokeTokensByUserId(userId uuid.UUID) error {
	// Xóa tất cả session của người dùng
	if err := global.Pdb.Where("user_id = ?", userId).Delete(&models.Session{}).Error; err != nil {
		return errors.New("failed to revoke tokens")
	}
	return nil
}
