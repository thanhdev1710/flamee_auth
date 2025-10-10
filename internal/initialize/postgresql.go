package initialize

import (
	"fmt"
	"log"
	"time"

	"github.com/thanhdev1710/flamee_auth/global"
	"github.com/thanhdev1710/flamee_auth/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitPostgreSql() {
	p := global.Config.Postgre

	var logLevel logger.LogLevel
	if global.Config.Env == "production" {
		logLevel = logger.Silent
	} else {
		logLevel = logger.Info
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		p.Host, p.Username, p.Password, p.Database, p.Port, p.SslMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: false,
		Logger:                 logger.Default.LogMode(logLevel),
	})
	if err != nil {
		panic(err)
	}

	global.Pdb = db

	if global.Pdb != nil {
		log.Println("✅ Kết nối PostgreSql thành công")
	}

	SetPool()

	// MigrateTables()
}

func SetPool() {
	p := global.Config.Postgre
	sqlDb, err := global.Pdb.DB()
	if err != nil {
		panic("postgreSql error")
	}
	sqlDb.SetConnMaxIdleTime(time.Duration(p.ConnMaxIdleTime))
	sqlDb.SetConnMaxLifetime(time.Duration(p.ConnMaxLifeTime))
	sqlDb.SetMaxOpenConns(p.ConnMaxOpen)
}

func MigrateTables() {
	if err := global.Pdb.AutoMigrate(
		&models.User{},
		&models.Session{},
		&models.OAuthProvider{},
		&models.VerificationToken{},
	); err != nil {
		panic("failed to migrate database: " + err.Error())
	}
}

func ClosePostgreSql() {
	// Lấy kết nối cơ sở dữ liệu
	sqlDb, err := global.Pdb.DB()
	if err != nil {
		panic("Lỗi khi lấy đối tượng DB PostgreSQL: " + err.Error())
	}

	// Đóng kết nối cơ sở dữ liệu
	err = sqlDb.Close()
	if err != nil {
		panic("Lỗi khi đóng kết nối PostgreSQL: " + err.Error())
	}

	fmt.Println("🔒 Đã đóng kết nối PostgreSQL")
}
