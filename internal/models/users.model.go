package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id                 uuid.UUID           `gorm:"type:uuid;primary_key" json:"id"`
	Email              string              `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password           string              `gorm:"type:varchar(255);not null" json:"password"`
	Role               string              `gorm:"type:varchar(10);not null" json:"role"`
	IsVerified         bool                `gorm:"default:false" json:"is_verified"`
	Status             string              `gorm:"type:varchar(10);default:inactive" json:"status"`
	CreatedAt          time.Time           `gorm:"default:now()" json:"created_at"`
	UpdatedAt          time.Time           `gorm:"default:now()" json:"updated_at"`
	Sessions           []Session           `gorm:"foreignKey:UserId" json:"sessions"`
	OAuthProviders     []OAuthProvider     `gorm:"foreignKey:UserId" json:"oauth_providers"`
	VerificationTokens []VerificationToken `gorm:"foreignKey:UserId" json:"verification_tokens"`
}
