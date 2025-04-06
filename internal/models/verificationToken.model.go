package models

import (
	"time"

	"github.com/google/uuid"
)

type VerificationToken struct {
	Id        int       `gorm:"primary_key"`
	UserId    uuid.UUID `gorm:"type:uuid;not null"`
	Token     string    `gorm:"type:varchar(255);unique;not null"`
	TokenType string    `gorm:"type:varchar(50);not null"` // email_verification, password_reset
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:now()"`
	User      User      `gorm:"foreignKey:UserId"` // Link to User model
}
