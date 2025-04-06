package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Id        int       `gorm:"primary_key"`
	UserId    uuid.UUID `gorm:"type:uuid;not null"`
	Token     string    `gorm:"type:varchar(512);unique;not null"`
	UserAgent string    `gorm:"type:text"`
	IpAddress string    `gorm:"type:varchar(45)"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:now()"`
	User      User      `gorm:"foreignKey:UserId"` // Link to User model
}
