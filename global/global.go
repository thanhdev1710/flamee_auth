package global

import (
	"github.com/thanhdev1710/flamee_auth/pkg/logger"
	"github.com/thanhdev1710/flamee_auth/pkg/settings"
	"gorm.io/gorm"
)

var (
	Config settings.Config
	Pdb    *gorm.DB
	Logger *logger.LoggerZap
)
