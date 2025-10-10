package initialize

import (
	"log"

	"github.com/spf13/viper"
	"github.com/thanhdev1710/flamee_auth/global"
)

func LoadConfig() {
	// Bind tất cả biến môi trường từ OS
	for _, key := range []string{
		"PORT", "API_KEY", "ENV", "DOMAIN",
		"URL_FRONT_END", "URL_POST_SERVICE", "URL_USER_SERVICE",
		"JWT_SECRET", "JWT_EXPIRATION_TIME_DEFAULT", "JWT_EXPIRATION_TIME_REMEMBER",
		"EMAIL_FROM", "EMAIL_PASSWORD", "EMAIL_SMTPHOST",
		"LOG_LEVEL", "LOG_FILE", "LOG_MAXSIZE", "LOG_MAXBACKUPS", "LOG_MAXAGE", "LOG_COMPRESS",
		"EMAIL_SMTPPORT", "EMAIL_SECRET",
		"HOST_DB", "PORT_DB", "USERNAME_DB", "PASSWORD_DB", "DATABASE_DB", "SSL_MODE_DB",
		"CONN_MAX_IDLE_TIME_DB", "CONN_MAX_OPEN", "CONN_MAX_LIFE_TIME",
		"NATS_HOST", "NATS_PORT", "NATS_MAX_RECONNECTS", "NATS_RECONNECT_WAIT",
	} {
		viper.BindEnv(key)
	}

	// Unmarshal vào struct config
	if err := viper.Unmarshal(&global.Config); err != nil {
		log.Fatalf("Cannot unmarshal config: %v", err)
	}
	if err := viper.Unmarshal(&global.Config.Postgre); err != nil {
		log.Fatalf("Cannot unmarshal config: %v", err)
	}
	if err := viper.Unmarshal(&global.Config.Email); err != nil {
		log.Fatalf("Cannot unmarshal config: %v", err)
	}
	if err := viper.Unmarshal(&global.Config.Logger); err != nil {
		log.Fatalf("Cannot unmarshal config: %v", err)
	}
	if err := viper.Unmarshal(&global.Url); err != nil {
		log.Fatalf("Cannot unmarshal config: %v", err)
	}
	if err := viper.Unmarshal(&global.Config.Nats); err != nil {
		log.Fatalf("Cannot unmarshal NATS config: %v", err)
	}

	log.Println("Config loaded from OS ENV")
}
