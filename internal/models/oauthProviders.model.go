package models

import (
	"time"

	"github.com/google/uuid"
)

type OAuthProvider struct {
	Id         int       `gorm:"primary_key"`
	UserId     uuid.UUID `gorm:"type:uuid;not null"`
	Provider   string    `gorm:"type:varchar(50);not null"` // google, facebook, github
	ProviderId string    `gorm:"type:varchar(255);unique;not null"`
	CreatedAt  time.Time `gorm:"default:now()"`
	User       User      `gorm:"foreignKey:UserId"` // Link to User model
}
