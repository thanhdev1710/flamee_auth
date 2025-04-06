package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id                 uuid.UUID           `gorm:"type:uuid;primary_key"`
	Email              string              `gorm:"type:varchar(100);unique;not null"`
	Password           string              `gorm:"type:varchar(255);not null"`
	Role               string              `gorm:"type:varchar(10);not null"`
	IsVerified         bool                `gorm:"default:false"`
	CreatedAt          time.Time           `gorm:"default:now()"`
	UpdatedAt          time.Time           `gorm:"default:now()"`
	Sessions           []Session           `gorm:"foreignKey:UserId"` // One-to-Many relationship with Session
	OAuthProviders     []OAuthProvider     `gorm:"foreignKey:UserId"` // One-to-Many relationship with OAuthProvider
	VerificationTokens []VerificationToken `gorm:"foreignKey:UserId"` // One-to-Many relationship with VerificationToken
}
