package global

import (
	"github.com/nats-io/nats.go"
	"github.com/thanhdev1710/flamee_auth/pkg/logger"
	"github.com/thanhdev1710/flamee_auth/pkg/settings"
	"gorm.io/gorm"
)

var (
	Config         settings.Config
	Pdb            *gorm.DB
	NatsConnection *nats.Conn
	Logger         *logger.LoggerZap
	User           = UserStatus{
		Active:   "active",
		Inactive: "inactive",
		Banned:   "banned",
		Deleted:  "deleted",
	}
	Url   settings.Url
	Token = tokenStruct{
		AccessToken:  "flamee_access_token",
		RefreshToken: "flamee_refresh_token",
	}
)

type tokenStruct struct {
	AccessToken  string
	RefreshToken string
}

type UserStatus struct {
	Active   string
	Inactive string
	Banned   string
	Deleted  string
}
